package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceGroupRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRulesRead,
		Schema: map[string]*schema.Schema{
			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": statusSchema,
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
						"users_excluded": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ruleName, nameOk := d.GetOk("name_prefix")
	searchParams := &query.Params{Limit: 200}
	if nameOk {
		searchParams.Search = ruleName.(string)
	}

	rules, resp, err := getOktaClientFromMetadata(m).Group.ListGroupRules(ctx, searchParams)
	if err != nil {
		return diag.Errorf("failed to list group rules: %v", err)
	}

	for {
		if resp.HasNextPage() {
			var nextRules []*sdk.GroupRule
			resp, err = resp.Next(ctx, &nextRules)
			if err != nil {
				return diag.Errorf("failed to list group rules: %v", err)
			}
			rules = append(rules, nextRules...)
		} else {
			break
		}
	}

	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(searchParams.String()))))
	rulesArr := make([]map[string]interface{}, len(rules))
	for i := range rules {
		rule := map[string]interface{}{
			"id":     rules[i].Id,
			"name":   rules[i].Name,
			"status": rules[i].Status,
		}

		if rules[i].Conditions != nil {
			rule["expression_type"] = rules[i].Conditions.Expression.Type
			rule["expression_value"] = rules[i].Conditions.Expression.Value
		}

		if rules[i].Conditions.People != nil && rules[i].Conditions.People.Users != nil {
			rule["users_excluded"] = convertStringSliceToSet(rules[i].Conditions.People.Users.Exclude)
		}

		rule["group_assignments"] = convertStringSliceToSet(rules[i].Actions.AssignUserToGroups.GroupIds)

		rulesArr[i] = rule
	}
	_ = d.Set("rules", rulesArr)
	return nil
}
