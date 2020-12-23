package okta

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

const statusInvalid = "INVALID"

func resourceGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupRuleCreate,
		ReadContext:   resourceGroupRuleRead,
		UpdateContext: resourceGroupRuleUpdate,
		DeleteContext: resourceGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_assignments": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// Actions cannot be updated even on a deactivated rule
				ForceNew: true,
			},
			"expression_type": {
				Type:     schema.TypeString,
				Default:  "urn:okta:expression:1.0",
				Optional: true,
			},
			"expression_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": statusSchema,
		},
		CustomizeDiff: customdiff.ForceNewIf("status", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			g, _, _ := getOktaClientFromMetadata(meta).Group.GetGroupRule(ctx, d.Id(), nil)
			if g == nil {
				return false
			}
			_ = d.SetNew("status", g.Status)
			return d.Get("status").(string) == statusInvalid
		}),
	}
}

func resourceGroupRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupRule := buildGroupRule(d)
	responseGroupRule, _, err := getOktaClientFromMetadata(m).Group.CreateGroupRule(ctx, *groupRule)
	if err != nil {
		return diag.Errorf("failed to create group rule: %v", err)
	}
	d.SetId(responseGroupRule.Id)
	if err := handleGroupRuleLifecycle(ctx, d, m); err != nil {
		return diag.Errorf("failed to change group rule status: %v", err)
	}
	return resourceGroupRuleRead(ctx, d, m)
}

func resourceGroupRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g, resp, err := getOktaClientFromMetadata(m).Group.GetGroupRule(ctx, d.Id(), nil)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get group rule: %v", err)
	}
	if g == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", g.Name)
	// _ = d.Set("type", g.Type)
	_ = d.Set("status", g.Status)
	// Just for the sake of safety, should never be nil
	if g.Conditions != nil && g.Conditions.Expression != nil {
		_ = d.Set("expression_type", g.Conditions.Expression.Type)
		_ = d.Set("expression_value", g.Conditions.Expression.Value)
	}
	err = setNonPrimitives(d, map[string]interface{}{
		"group_assignments": convertStringSetToInterface(g.Actions.AssignUserToGroups.GroupIds),
	})
	if err != nil {
		return diag.Errorf("failed to set group rule properties: %v", err)
	}
	return nil
}

func resourceGroupRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	desiredStatus := d.Get("status").(string)
	// Only inactive rules can be changed, thus we should handle this first
	if d.HasChange("status") {
		err := handleGroupRuleLifecycle(ctx, d, m)
		if err != nil {
			return diag.Errorf("failed to change group rule status: %v", err)
		}
		_ = d.Set("status", desiredStatus)
	}
	// invalid group rules can not be updated
	if hasGroupRuleChange(d) && desiredStatus != statusInvalid {
		client := getOktaClientFromMetadata(m)
		rule := buildGroupRule(d)
		if desiredStatus == statusActive {
			// Only inactive rules can be changed, thus we should deactivate the rule in case it was "ACTIVE"
			_, err := client.Group.DeactivateGroupRule(ctx, d.Id())
			if err != nil {
				return diag.Errorf("failed to deactivate group rule: %v", err)
			}
		}
		_, _, err := client.Group.UpdateGroupRule(ctx, d.Id(), *rule)
		if err != nil {
			return diag.Errorf("failed to update group rule: %v", err)
		}
		if desiredStatus == statusActive {
			// We should reactivate the rule in case it was deactivated.
			_, err := client.Group.ActivateGroupRule(ctx, d.Id())
			if err != nil {
				return diag.Errorf("failed to activate group rule: %v", err)
			}
		}
	}
	return resourceGroupRuleRead(ctx, d, m)
}

func hasGroupRuleChange(d *schema.ResourceData) bool {
	for _, k := range []string{"expression_type", "expression_value", "name", "group_assignments"} {
		if d.HasChange(k) {
			return true
		}
	}
	return false
}

func resourceGroupRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	if d.Get("status").(string) == statusActive {
		_, err := client.Group.DeactivateGroupRule(ctx, d.Id())
		// suppress error for INACTIVE group rules
		if err != nil && !strings.Contains(err.Error(), "Cannot activate or deactivate a Group Rule with the status INVALID") {
			return diag.Errorf("failed to deactivate group rule before removing: %v", err)
		}
	}
	_, err := client.Group.DeleteGroupRule(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete group rule: %v", err)
	}
	return nil
}

func buildGroupRule(d *schema.ResourceData) *okta.GroupRule {
	return &okta.GroupRule{
		Actions: &okta.GroupRuleAction{
			AssignUserToGroups: &okta.GroupRuleGroupAssignment{
				GroupIds: convertInterfaceToStringSet(d.Get("group_assignments")),
			},
		},
		Conditions: &okta.GroupRuleConditions{
			Expression: &okta.GroupRuleExpression{
				Type:  d.Get("expression_type").(string),
				Value: d.Get("expression_value").(string),
			},
		},
		Name: d.Get("name").(string),
		Type: "group_rule",
	}
}

func handleGroupRuleLifecycle(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	if d.Get("status").(string) == statusActive {
		_, err := client.Group.ActivateGroupRule(ctx, d.Id())
		return err
	} else if d.Get("status").(string) == statusInvalid {
		return nil
	}
	_, err := client.Group.DeactivateGroupRule(ctx, d.Id())
	return err
}
