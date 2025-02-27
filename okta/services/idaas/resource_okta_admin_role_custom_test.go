package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaAdminRoleCustom_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAdminRoleCustom, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAdminRoleCustom)
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAdminRoleCustom, doesAdminRoleCustomExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing"),
						resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing updated"),
						resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
					),
				},
			},
		})
}

func doesAdminRoleCustomExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
	_, response, err := client.GetCustomRole(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}
