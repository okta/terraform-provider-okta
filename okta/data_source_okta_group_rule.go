package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceGroupRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_assignments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"expression_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expression_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": statusSchema,
			"users_excluded": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGroupRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var rule *sdk.GroupRule
	ruleID, idOk := d.GetOk("id")
	if idOk {
		respRule, _, err := getOktaClientFromMetadata(m).Group.GetGroupRule(ctx, ruleID.(string), nil)
		if err != nil {
			return diag.Errorf("failed get group rule by ID: %v", err)
		}
		rule = respRule
	} else {
		ruleName, nameOk := d.GetOk("name")
		if nameOk {
			name := ruleName.(string)
			searchParams := &query.Params{Search: name, Limit: 1}
			rules, _, err := getOktaClientFromMetadata(m).Group.ListGroupRules(ctx, searchParams)
			switch {
			case err != nil:
				return diag.Errorf("failed to get group rule by name: %v", err)
			case len(rules) < 1:
				return diag.Errorf("group rule with name '%s' does not exist", name)
			}
			// exact name search should only return one result, but loop through to be safe
			for _, ruleCandidate := range rules {
				if ruleName == ruleCandidate.Name {
					rule = ruleCandidate
					break
				}
			}
		}
	}

	if rule == nil {
		return diag.Errorf("config must provide 'name' or 'id' to retrieve a group rule")
	}

	d.SetId(rule.Id)
	_ = d.Set("name", rule.Name)
	_ = d.Set("status", rule.Status)
	if rule.Conditions != nil {
		_ = d.Set("expression_type", rule.Conditions.Expression.Type)
		_ = d.Set("expression_value", rule.Conditions.Expression.Value)
	}
	if rule.Conditions.People != nil && rule.Conditions.People.Users != nil {
		_ = d.Set("users_excluded", convertStringSliceToSet(rule.Conditions.People.Users.Exclude))
	}
	err := setNonPrimitives(d, map[string]interface{}{
		"group_assignments": convertStringSliceToSet(rule.Actions.AssignUserToGroups.GroupIds),
	})
	if err != nil {
		return diag.Errorf("failed to set group rule properties: %v", err)
	}
	return nil
}
