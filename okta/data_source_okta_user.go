package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: buildUserDataSourceSchema(map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Description: "Retrieve a single user based on their id",
				Optional:    true,
			},
			"search": {
				Type:        schema.TypeSet,
				Optional:    true,
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
		}),
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	var user *okta.User
	var err error
	userID, ok := d.GetOk("user_id")
	_, searchCriteriaOk := d.GetOk("search")
	if ok {
		logger(m).Info("reading user by ID", "id", userID.(string))
		user, _, err = client.User.GetUser(ctx, userID.(string))
		if err != nil {
			return diag.Errorf("failed to get user: %v", err)
		}
	} else if searchCriteriaOk {
		var users []*okta.User
		sc := getSearchCriteria(d)
		logger(m).Info("reading user using search", "search", sc)
		users, _, err = client.User.ListUsers(ctx, &query.Params{Search: sc, Limit: 1})
		if err != nil {
			return diag.Errorf("failed to list users: %v", err)
		} else if len(users) < 1 {
			return diag.Errorf("no users found using search criteria: %+v", sc)
		}
		user = users[0]
	}

	d.SetId(user.Id)

	rawMap := flattenUser(user)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set user's properties: %v", err)
	}
	err = setAdminRoles(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to set user's admin roles: %v", err)
	}
	return nil
}

func getSearchCriteria(d *schema.ResourceData) string {
	rawFilters := d.Get("search").(*schema.Set)
	filterList := make([]string, rawFilters.Len())
	for i, f := range rawFilters.List() {
		fmap := f.(map[string]interface{})
		filterList[i] = fmt.Sprintf(`%s %s "%s"`, fmap["name"], fmap["comparison"], fmap["value"])
	}
	return strings.Join(filterList, " and ")
}
