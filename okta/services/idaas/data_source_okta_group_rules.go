package idaas

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceGroupRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRulesRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Searches for group rules with a supported filtering expression for all attributes except for '_embedded', '_links', and 'objectClass'",
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The maximum number of group rules returned by the Okta API, between 1 and 200.",
				Default:      utils.DefaultPaginationLimit,
				ValidateFunc: validation.IntBetween(1, 200),
			},
			"expand": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If specified as `groupIdToGroupNameMap`, then displays group names",
			},
			"group_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group rule ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group rule name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group rule status.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group rule type.",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group rule creation date.",
						},
						"last_updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group rule last updated date.",
						},
						"expression_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expression type to use to invoke the rule.",
						},
						"expression_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expression value.",
						},
						"group_assignments": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of group ids to assign the users to.",
						},
						"users_excluded": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of user IDs that would be excluded when rules are processed.",
						},
					},
				},
			},
		},
		Description: "Get a list of group rules from Okta.",
	}
}

func dataSourceGroupRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)

	qp := &query.Params{Limit: utils.DefaultPaginationLimit}

	if limit, ok := d.GetOk("limit"); ok {
		qp.Limit = int64(limit.(int))
	}

	if search, ok := d.GetOk("search"); ok {
		qp.Search = search.(string)
	}

	if expand, ok := d.GetOk("expand"); ok {
		qp.Expand = expand.(string)
	}

	groupRules, resp, err := client.Group.ListGroupRules(ctx, qp)
	if err != nil {
		d.SetId("")
		return diag.Errorf("failed to list group rules: %v", err)
	}

	// handle pagination
	for {
		if !resp.HasNextPage() {
			break
		}
		var moreRules []*sdk.GroupRule
		var err error
		resp, err = resp.Next(ctx, &moreRules)
		if err != nil {
			return diag.Errorf("failed to get next page of group rules: %v", err)
		}
		groupRules = append(groupRules, moreRules...)
	}

	// generate a unique ID for the data source based on the query parameters
	dataSourceId := fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(qp.String())))
	d.SetId(dataSourceId)

	// convert the group rules to a list of maps
	arr := make([]map[string]interface{}, len(groupRules))
	for i := range groupRules {
		arr[i] = map[string]interface{}{}

		arr[i]["id"] = groupRules[i].Id
		arr[i]["name"] = groupRules[i].Name
		arr[i]["status"] = groupRules[i].Status
		arr[i]["type"] = groupRules[i].Type

		if groupRules[i].Created != nil {
			arr[i]["created"] = groupRules[i].Created.Format("2006-01-02T15:04:05.000Z")
		}
		if groupRules[i].LastUpdated != nil {
			arr[i]["last_updated"] = groupRules[i].LastUpdated.Format("2006-01-02T15:04:05.000Z")
		}

		if groupRules[i].Conditions != nil && groupRules[i].Conditions.Expression != nil {
			arr[i]["expression_type"] = groupRules[i].Conditions.Expression.Type
			arr[i]["expression_value"] = groupRules[i].Conditions.Expression.Value
		}

		if groupRules[i].Actions != nil && groupRules[i].Actions.AssignUserToGroups != nil {
			arr[i]["group_assignments"] = utils.ConvertStringSliceToSet(groupRules[i].Actions.AssignUserToGroups.GroupIds)
		}

		if groupRules[i].Conditions != nil && groupRules[i].Conditions.People != nil && groupRules[i].Conditions.People.Users != nil {
			arr[i]["users_excluded"] = utils.ConvertStringSliceToSet(groupRules[i].Conditions.People.Users.Exclude)
		}
	}

	err = d.Set("group_rules", arr)
	if err != nil {
		return diag.Errorf("failed to set group rules: %v", err)
	}

	return nil
} 
