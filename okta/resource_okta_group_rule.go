package okta

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

const statusInvalid = "INVALID"
const maxRetries = 3
const retryDelay = 60 * time.Second

func resourceGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupRuleCreate,
		ReadContext:   resourceGroupRuleRead,
		UpdateContext: resourceGroupRuleUpdate,
		DeleteContext: resourceGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Creates an Okta Group Rule.
This resource allows you to create and configure an Okta Group Rule.
-> If the Okta API marks the 'status' of the rule as 'INVALID' the Okta
Terraform Provider will act in a force/replace manner and call the API to delete
the underlying rule resource and create a new rule resource.`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The name of the Group Rule (min character 1; max characters 50).",
				ValidateDiagFunc: strMaxLength(50), // ideas.okta.com/app/#/case/204022
			},
			"group_assignments": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// Actions cannot be updated even on a deactivated rule
				ForceNew:    true,
				Description: "The list of group ids to assign the users to.",
			},
			"expression_type": {
				Type:        schema.TypeString,
				Default:     "urn:okta:expression:1.0",
				Optional:    true,
				Description: "The expression type to use to invoke the rule. The default is `urn:okta:expression:1.0`.",
			},
			"expression_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The expression value.",
			},
			"status": statusSchema,
			"remove_assigned_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Remove users added by this rule from the assigned group after deleting this resource. Default is `false`",
			},
			"users_excluded": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The list of user IDs that would be excluded when rules are processed",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		CustomizeDiff: customdiff.ForceNewIf("status", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			return statusIsInvalidDiffFn(d.Get("status").(string))
		}),
	}
}

func statusIsInvalidDiffFn(status string) bool {
	return status == statusInvalid
}

func resourceGroupRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Using the existing retryOnStatusCodes context value for INTERNAL SERVER ERROR
	ctx = context.WithValue(ctx, retryOnStatusCodes, []int{http.StatusInternalServerError})
	groupRule := buildGroupRule(d)
	client := getOktaClientFromMetadata(meta)

	var responseGroupRule *sdk.GroupRule
	var sdkResp *sdk.Response
	var err error

	// Retry loop for CreateGroupRule in case of HTTP 423 responses
	for attempt := 0; attempt < maxRetries; attempt++ {
		responseGroupRule, sdkResp, err = client.Group.CreateGroupRule(ctx, *groupRule)
		if err != nil && sdkResp != nil && sdkResp.StatusCode == http.StatusLocked {
			time.Sleep(retryDelay)
			continue
		}
		break
	}
	if err != nil {
		return diag.Errorf("failed to create group rule: %v", err)
	}
	d.SetId(responseGroupRule.Id)
	// Retry lifecycle operations as needed (e.g. activation or deactivation)
	if err := handleGroupRuleLifecycle(ctx, d, meta); err != nil {
		return diag.Errorf("failed to change group rule status: %v", err)
	}
	return resourceGroupRuleRead(ctx, d, meta)
}

func resourceGroupRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	g, resp, err := getOktaClientFromMetadata(meta).Group.GetGroupRule(ctx, d.Id(), nil)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get group rule: %v", err)
	}
	if g == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", g.Name)
	_ = d.Set("status", g.Status)
	// Just for the sake of safety, should never be nil
	if g.Conditions != nil {
		if g.Conditions.Expression != nil {
			_ = d.Set("expression_type", g.Conditions.Expression.Type)
			_ = d.Set("expression_value", g.Conditions.Expression.Value)
		}
		if g.Conditions.People != nil && g.Conditions.People.Users != nil {
			_ = d.Set("users_excluded", convertStringSliceToSet(g.Conditions.People.Users.Exclude))
		}
	}
	err = setNonPrimitives(d, map[string]interface{}{
		"group_assignments": convertStringSliceToSet(g.Actions.AssignUserToGroups.GroupIds),
	})
	if err != nil {
		return diag.Errorf("failed to set group rule properties: %v", err)
	}
	return nil
}

func resourceGroupRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	desiredStatus := d.Get("status").(string)
	client := getOktaClientFromMetadata(meta)
	// Only inactive rules can be changed, thus we should handle this first
	if d.HasChange("status") {
		if err := handleGroupRuleLifecycle(ctx, d, meta); err != nil {
			return diag.Errorf("failed to change group rule status: %v", err)
		}
		_ = d.Set("status", desiredStatus)
	}
	// If any properties (other than status) have changed and the desired status is not INVALID,
	// we need to update the rule. We must deactivate first if the rule is active.
	if hasGroupRuleChange(d) && desiredStatus != statusInvalid {
		rule := buildGroupRule(d)
		// If rule is active, deactivate before updating
		if desiredStatus == statusActive {
			for attempt := 0; attempt < maxRetries; attempt++ {
				_, err := client.Group.DeactivateGroupRule(ctx, d.Id())
				if err != nil && strings.Contains(err.Error(), "423") {
					time.Sleep(retryDelay)
					continue
				}
				if err != nil {
					return diag.Errorf("failed to deactivate group rule: %v", err)
				}
				break
			}
		}
		// Retry loop for UpdateGroupRule
		var sdkResp *sdk.Response
		var err error
		for attempt := 0; attempt < maxRetries; attempt++ {
			_, sdkResp, err = client.Group.UpdateGroupRule(ctx, d.Id(), *rule)
			if err != nil && sdkResp != nil && sdkResp.StatusCode == http.StatusLocked {
				time.Sleep(retryDelay)
				continue
			}
			break
		}
		if err != nil {
			return diag.Errorf("failed to update group rule: %v", err)
		}
		if desiredStatus == statusActive {
			// We should reactivate the rule in case it was deactivated.
			for attempt := 0; attempt < maxRetries; attempt++ {
				_, err := client.Group.ActivateGroupRule(ctx, d.Id())
				if err != nil && strings.Contains(err.Error(), "423") {
					time.Sleep(retryDelay)
					continue
				}
				if err != nil {
					return diag.Errorf("failed to activate group rule: %v", err)
				}
				break
			}
		}
	}
	return resourceGroupRuleRead(ctx, d, meta)
}

func hasGroupRuleChange(d *schema.ResourceData) bool {
	for _, k := range []string{"expression_type", "expression_value", "name", "group_assignments", "users_excluded"} {
		if d.HasChange(k) {
			return true
		}
	}
	return false
}

func resourceGroupRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)
	// If the rule is active, attempt to deactivate it first.
	if d.Get("status").(string) == statusActive {
		for attempt := 0; attempt < maxRetries; attempt++ {
			_, err := client.Group.DeactivateGroupRule(ctx, d.Id())
			// If error is due to a 423 response, retry.
			if err != nil && strings.Contains(err.Error(), "423") {
				time.Sleep(retryDelay)
				continue
			}
			// suppress error for INACTIVE group rules
			if err != nil && strings.Contains(err.Error(), "Cannot activate or deactivate a Group Rule with the status INVALID") {
				err = nil
			}
			if err != nil {
				return diag.Errorf("failed to deactivate group rule before removing: %v", err)
			}
			break
		}
	}
	remove := d.Get("remove_assigned_users").(bool)
	var err error
	// Retry loop for DeleteGroupRule
	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err = client.Group.DeleteGroupRule(ctx, d.Id(), &query.Params{RemoveUsers: &remove})
		if err != nil && strings.Contains(err.Error(), "423") {
			time.Sleep(retryDelay)
			continue
		}
		break
	}
	if err != nil {
		return diag.Errorf("failed to delete group rule: %v", err)
	}
	return nil
}

func buildGroupRule(d *schema.ResourceData) *sdk.GroupRule {
	return &sdk.GroupRule{
		Actions: &sdk.GroupRuleAction{
			AssignUserToGroups: &sdk.GroupRuleGroupAssignment{
				GroupIds: convertInterfaceToStringSet(d.Get("group_assignments")),
			},
		},
		Conditions: &sdk.GroupRuleConditions{
			Expression: &sdk.GroupRuleExpression{
				Type:  d.Get("expression_type").(string),
				Value: d.Get("expression_value").(string),
			},
			People: &sdk.GroupRulePeopleCondition{
				Users: &sdk.GroupRuleUserCondition{
					Exclude: convertInterfaceToStringSet(d.Get("users_excluded")),
				},
			},
		},
		Name: d.Get("name").(string),
		Type: "group_rule",
	}
}

func handleGroupRuleLifecycle(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := getOktaClientFromMetadata(meta)
	if d.Get("status").(string) == statusActive {
		for attempt := 0; attempt < maxRetries; attempt++ {
			_, err := client.Group.ActivateGroupRule(ctx, d.Id())
			if err != nil && strings.Contains(err.Error(), "423") {
				time.Sleep(retryDelay)
				continue
			}
			return err
		}
	} else if d.Get("status").(string) == statusInvalid {
		return nil
	}
	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err := client.Group.DeactivateGroupRule(ctx, d.Id())
		if err != nil && strings.Contains(err.Error(), "423") {
			time.Sleep(retryDelay)
			continue
		}
		return err
	}
	return nil
}
