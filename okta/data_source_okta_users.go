package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
					Schema: userProfileDataSchema,
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	results := &searchResults{Users: []*okta.User{}}
	params := &query.Params{Search: getSearchCriteria(d), Limit: defaultPaginationLimit, SortOrder: "0"}
	err := collectUsers(ctx, client, results, params)
	if err != nil {
		return diag.Errorf("failed to list users: %v", err)
	}

	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(params.String()))))
	arr := make([]map[string]interface{}, len(results.Users))

	for i, user := range results.Users {
		rawMap := flattenUser(user)
		arr[i] = rawMap
	}
	_ = d.Set("users", arr)

	return nil
}

// Recursively list apps until no next links are returned
func collectUsers(ctx context.Context, client *okta.Client, results *searchResults, qp *query.Params) error {
	users, res, err := client.User.ListUsers(ctx, qp)
	if err != nil {
		return err
	}

	results.Users = append(results.Users, users...)

	if after := sdk.GetAfterParam(res); after != "" {
		qp.After = after
		return collectUsers(ctx, client, results, qp)
	}

	return nil
}
