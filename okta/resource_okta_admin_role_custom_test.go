package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaAdminRoleCustom_crud(t *testing.T) {
	mgr := newFixtureManager("resources", adminRoleCustom, t.Name())
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
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
}

func TestAccResourceOktaAdminRoleCustom_withConditions(t *testing.T) {
	mgr := newFixtureManager("resources", adminRoleCustom, t.Name())
	config := mgr.GetFixtures("basic_with_conditions.tf", t)
	resourceName := fmt.Sprintf("%s.test", adminRoleCustom)
	resourceName1 := fmt.Sprintf("%s.test1", adminRoleCustom)
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
						resource.TestCheckResourceAttr(resourceName, "description", "Testing custom role with permission conditions"),
						resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "permission_conditions.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "permission_conditions.0.permission", "okta.users.read"),
						resource.TestCheckResourceAttr(resourceName, "permission_conditions.0.include", `{"okta:ResourceAttribute/User/Profile":["department","costCenter"]}`),

						resource.TestCheckResourceAttr(resourceName1, "label", buildResourceName(mgr.Seed)+"_1"),
						resource.TestCheckResourceAttr(resourceName1, "description", "Testing custom role with permission conditions"),
						resource.TestCheckResourceAttr(resourceName1, "permissions.#", "1"),
						resource.TestCheckResourceAttr(resourceName1, "permission_conditions.#", "1"),
						resource.TestCheckResourceAttr(resourceName1, "permission_conditions.0.permission", "okta.users.read"),
						resource.TestCheckResourceAttr(resourceName1, "permission_conditions.0.exclude", `{"okta:ResourceAttribute/User/Profile":["title"]}`),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      resourceName1,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
}

func doesAdminRoleCustomExist(id string) (bool, error) {
	client := sdkSupplementClientForTest()
	_, response, err := client.GetCustomRole(context.Background(), id)
	return doesResourceExist(response, err)
}
