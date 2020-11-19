package okta

import (
	"context"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupRuleCreate,
		Exists: resourceGroupRuleExists,
		Read:   resourceGroupRuleRead,
		Update: resourceGroupRuleUpdate,
		Delete: resourceGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
	}
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

func handleGroupRuleLifecycle(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)

	if d.Get("status").(string) == statusActive {
		_, err := client.Group.ActivateGroupRule(context.Background(), d.Id())
		return err
	}

	_, err := client.Group.DeactivateGroupRule(context.Background(), d.Id())
	return err
}

func resourceGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	groupRule := buildGroupRule(d)
	responseGroupRule, _, err := getOktaClientFromMetadata(m).Group.CreateGroupRule(context.Background(), *groupRule)
	if err != nil {
		return err
	}
	d.SetId(responseGroupRule.Id)

	if err := handleGroupRuleLifecycle(d, m); err != nil {
		return err
	}

	return resourceGroupRuleRead(d, m)
}

func resourceGroupRuleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchGroupRule(d, m)

	return err == nil && g != nil, err
}

func resourceGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	g, err := fetchGroupRule(d, m)

	if g == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("name", g.Name)
	_ = d.Set("type", g.Type)
	_ = d.Set("status", g.Status)

	// Just for the sake of safety, should never be nil
	if g.Conditions != nil && g.Conditions.Expression != nil {
		_ = d.Set("expression_type", g.Conditions.Expression.Type)
		_ = d.Set("expression_value", g.Conditions.Expression.Value)
	}

	return setNonPrimitives(d, map[string]interface{}{
		"group_assignments": convertStringSetToInterface(g.Actions.AssignUserToGroups.GroupIds),
	})
}

func resourceGroupRuleUpdate(d *schema.ResourceData, m interface{}) error {
	desiredStatus := d.Get("status").(string)
	// Only inactive rules can be changed, thus we should handle this first
	if d.HasChange("status") {
		if err := handleGroupRuleLifecycle(d, m); err != nil {
			return err
		}
		_ = d.Set("status", desiredStatus)
	}

	if hasGroupRuleChange(d) {
		client := getOktaClientFromMetadata(m)
		rule := buildGroupRule(d)

		if desiredStatus == statusActive {
			// Only inactive rules can be changed, thus we should deactivate the rule in case it was "ACTIVE"
			if _, err := client.Group.DeactivateGroupRule(context.Background(), d.Id()); err != nil {
				return err
			}
		}

		_, _, err := client.Group.UpdateGroupRule(context.Background(), d.Id(), *rule)
		if err != nil {
			return err
		}

		if desiredStatus == statusActive {
			// We should reactivate the rule in case it was deactivated.
			if _, err := client.Group.ActivateGroupRule(context.Background(), d.Id()); err != nil {
				return err
			}
		}
	}

	return resourceGroupRuleRead(d, m)
}

func hasGroupRuleChange(d *schema.ResourceData) bool {
	for _, k := range []string{"expression_type", "expression_value", "name", "group_assignments"} {
		if d.HasChange(k) {
			return true
		}
	}
	return false
}

func resourceGroupRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	if _, err := client.Group.DeactivateGroupRule(context.Background(), d.Id()); err != nil {
		return err
	}

	_, err := client.Group.DeleteGroupRule(context.Background(), d.Id())

	return err
}

func fetchGroupRule(d *schema.ResourceData, m interface{}) (*okta.GroupRule, error) {
	g, resp, err := getOktaClientFromMetadata(m).Group.GetGroupRule(context.Background(), d.Id(), nil)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return g, err
}
