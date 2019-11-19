package okta

import (
	"fmt"

	"github.com/articulate/terraform-provider-okta/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUsersRead,

		Schema: map[string]*schema.Schema{
			"search": userSearchSchema,
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: userProfileDataSchema,
				},
			},
		},
	}
}

func dataSourceUsersRead(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	results := &searchResults{Users: []*okta.User{}}
	params := &query.Params{Search: getSearchCriteria(d), Limit: 200, SortOrder: "0"}
	err := collectUsers(client, results, params)
	if err != nil {
		return fmt.Errorf("Error Getting User from Okta: %v", err)
	}

	d.SetId(fmt.Sprintf("%d", hashcode.String(params.String())))
	arr := make([]map[string]interface{}, len(results.Users))

	for i, user := range results.Users {
		rawMap, err := flattenUser(user, d)
		if err != nil {
			return err
		}
		arr[i] = rawMap
	}

	return d.Set("users", arr)
}

// Recursively list apps until no next links are returned
func collectUsers(client *okta.Client, results *searchResults, qp *query.Params) error {
	users, res, err := client.User.ListUsers(qp)
	if err != nil {
		return err
	}

	results.Users = append(results.Users, users...)

	if after := sdk.GetAfterParam(res); after != "" {
		qp.After = after
		return collectUsers(client, results, qp)
	}

	return nil
}
