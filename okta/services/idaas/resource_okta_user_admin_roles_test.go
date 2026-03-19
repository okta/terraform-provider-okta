package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaUserAdminRoles_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUserAdminRoles, t.Name())
	start := mgr.GetFixtures("basic.tf", t)
	update := mgr.GetFixtures("basic_update.tf", t)
	remove := mgr.GetFixtures("basic_removal.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUserAdminRoles)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: start,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "2"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "3"),
				),
			},
			{
				Config: remove,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "1"),
				),
			},
		},
	})
}
