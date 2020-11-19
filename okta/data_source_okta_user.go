package okta

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

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
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "eq",
							ValidateFunc: validation.StringInSlice([]string{"eq", "lt", "gt", "sw"}, true),
						},
					},
				},
			},
		}),
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)

	var user *okta.User
	var err error

	userID, ok := d.GetOk("user_id")
	_, searchCriteriaOk := d.GetOk("search")

	if ok {
		user, _, err = client.User.GetUser(context.Background(), userID.(string))
		if err != nil {
			return err
		}
	} else if searchCriteriaOk {
		var users []*okta.User
		users, _, err = client.User.ListUsers(context.Background(), &query.Params{Search: getSearchCriteria(d), Limit: 1})

		if err != nil {
			return err
		} else if len(users) < 1 {
			return errors.New("failed to locate user with provided parameters")
		}

		user = users[0]
	}

	d.SetId(user.Id)

	rawMap, err := flattenUser(user)
	if err != nil {
		return err
	}

	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return err
	}

	return setAdminRoles(d, client)
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
