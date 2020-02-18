package okta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"search": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter to find a user, each filter will be concatenated with an AND clause. Please be aware profile properties must match what is in Okta, which is likely camel case",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property name to search for. This requires the search feature be on. Please see Okta documentation on their filter API for users. https://developer.okta.com/docs/api/resources/users#list-users-with-search",
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"comparison": &schema.Schema{
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

	userId, userIdOk := d.GetOk("user_id")
	_, searchCriteriaOk := d.GetOk("search")

	if userIdOk {
		user, _, err = client.User.GetUser(userId.(string))
		if err != nil {
			return err
		}
	} else if searchCriteriaOk {
		var users []*okta.User
		users, _, err := client.User.ListUsers(&query.Params{Search: getSearchCriteria(d), Limit: 1})

		if err != nil {
			return err
		} else if len(users) < 1 {
			return errors.New("failed to locate user with provided parameters")
		}

		user = users[0]
	}

	d.SetId(user.Id)

	rawMap, err := flattenUser(user, d)
	if err != nil {
		return err
	}

	if err = setNonPrimitives(d, rawMap); err != nil {
		return err
	}

	if err = setAdminRoles(d, client); err != nil {
		return err
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
