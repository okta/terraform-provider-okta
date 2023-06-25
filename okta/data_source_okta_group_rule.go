package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroupRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
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
	ruleID, ok := d.GetOk("id")
	if ok {
		rule, _, err := getOktaClientFromMetadata(m).Group.GetGroupRule(ctx, ruleID.(string), nil)
		if err != nil {
			return diag.Errorf("failed get group rule by ID: %v", err)
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
		err = setNonPrimitives(d, map[string]interface{}{
			"group_assignments": convertStringSliceToSet(rule.Actions.AssignUserToGroups.GroupIds),
		})
		if err != nil {
			return diag.Errorf("failed to set group rule properties: %v", err)
		}
		return nil
	} else {
		return diag.Errorf("config must provide 'id' to retrieve a group rule")
	}
}
