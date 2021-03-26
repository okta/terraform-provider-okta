package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Filter to find a user, each filter will be concatenated with an AND clause. Please be aware profile properties must match what is in Okta, which is likely camel case",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property name to search for. This requires the search feature be on. Please see Okta documentation on their filter API for users. https://developer.okta.com/docs/api/resources/users#list-users-with-search",
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comparison": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "eq",
							ValidateDiagFunc: stringInSlice([]string{"eq", "lt", "gt", "sw"}),
						},
					},
				},
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: buildSchema(userProfileDataSchema,
						map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeString,
								Computed: true,
							},
						}),
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := &query.Params{Search: getSearchCriteria(d), Limit: defaultPaginationLimit, SortOrder: "0"}
	users, err := collectUsers(ctx, getOktaClientFromMetadata(m), params)
	if err != nil {
		return diag.Errorf("failed to list users: %v", err)
	}
	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(params.String()))))
	arr := make([]map[string]interface{}, len(users))
	for i, user := range users {
		rawMap := flattenUser(user)
		rawMap["id"] = user.Id
		arr[i] = rawMap
	}
	_ = d.Set("users", arr)
	return nil
}

func collectUsers(ctx context.Context, client *okta.Client, qp *query.Params) ([]*okta.User, error) {
	users, resp, err := client.User.ListUsers(ctx, qp)
	if err != nil {
		return nil, err
	}
	for resp.HasNextPage() {
		var nextUsers []*okta.User
		resp, err = resp.Next(ctx, &nextUsers)
		if err != nil {
			return nil, err
		}
		for i := range nextUsers {
			users = append(users, nextUsers[i])
		}
	}
	return users, nil
}
