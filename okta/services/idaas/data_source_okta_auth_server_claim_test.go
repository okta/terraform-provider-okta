package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAuthServerClaim_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthServerClaim, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	createUser := mgr.GetFixtures("datasource_create_auth_server.tf", t)
	resourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSAuthServerClaim)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: createUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_auth_server.test", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "IDENTITY"),
				),
			},
		},
	})
}
