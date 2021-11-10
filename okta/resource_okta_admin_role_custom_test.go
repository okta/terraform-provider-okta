package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func sweepCustomRoles(client *testClient) error {
	var errorList []error
	customRoles, _, err := client.apiSupplement.ListCustomRoles(context.Background(), &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return err
	}
	for _, role := range customRoles {
		if !strings.HasPrefix(role.Label, "testAcc_") {
			if _, err := client.apiSupplement.DeleteCustomRole(context.Background(), role.Id); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func TestAccOktaAdminRoleCustom(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(adminRoleCustom)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", adminRoleCustom)
	resource.Test(
		t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createCheckResourceDestroy(adminRoleCustom, doesAdminRoleCustomExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing"),
						resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing updated"),
						resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
					),
				},
			},
		})
}

func doesAdminRoleCustomExist(id string) (bool, error) {
	_, response, err := getSupplementFromMetadata(testAccProvider.Meta()).GetCustomRole(context.Background(), id)
	return doesResourceExist(response, err)
}
