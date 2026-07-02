package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAdminRoleCustom_read(t *testing.T) {
	datasourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSAdminRoleCustom)

	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAdminRoleCustom, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttr(datasourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(datasourceName, "description", "test custom role"),
					resource.TestCheckResourceAttr(datasourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "permissions.0", "okta.apps.assignment.manage"),
				),
			},
		},
	})
}
