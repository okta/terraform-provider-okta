package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaAdminRoleCustom(t *testing.T) {
	mgr := newFixtureManager(adminRoleCustom, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", adminRoleCustom)
	oktaResourceTest(
		t, resource.TestCase{
			PreCheck:          testAccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      checkResourceDestroy(adminRoleCustom, doesAdminRoleCustomExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing"),
						resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing updated"),
						resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
					),
				},
			},
		})
}

func doesAdminRoleCustomExist(id string) (bool, error) {
	client := sdkSupplementClientForTest()
	_, response, err := client.GetCustomRole(context.Background(), id)
	return doesResourceExist(response, err)
}
