package okta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/okta/okta-sdk-golang/okta/query"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema: buildUserDataSourceSchema(map[string]*schema.Schema{
			"search": userSearchSchema,
		}),
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	users, _, err := client.User.ListUsers(&query.Params{Search: getSearchCriteria(d), Limit: 1})

	if err != nil {
		return fmt.Errorf("Error Getting User from Okta: %v", err)
	} else if len(users) < 1 {
		return errors.New("failed to locate user with provided parameters")
	}

	d.SetId(users[0].Id)
	rawMap, err := flattenUser(users[0], d)
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
