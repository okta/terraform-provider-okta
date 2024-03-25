package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaAdminRoleCustom_read(t *testing.T) {
	mgr := newFixtureManager("datasources", adminRoleCustom, t.Name())
	roleCreate := adminRoleCustomResource
	roleRead := fmt.Sprintf("%s%s", adminRoleCustomResource, adminRoleCustomDataSources)
	resourceName := fmt.Sprintf("%s.test", adminRoleCustom)
	dataSourceName := fmt.Sprintf("data.%s.test", adminRoleCustom)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(roleCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "Test admin role"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
				),
			},
			{
				Config: mgr.ConfigReplace(roleRead),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(dataSourceName, "description", "Test admin role"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
				),
			},
		},
	})
}

const adminRoleCustomResource = `
resource "okta_admin_role_custom" "test" {
	label       = "testAcc_replace_with_uuid"
	description = "Test admin role"
	permissions = ["okta.users.read"]
}
`

const adminRoleCustomDataSources = `
data "okta_admin_role_custom" "test" {
	label = "testRole_replace_with_uuid"
}
`
