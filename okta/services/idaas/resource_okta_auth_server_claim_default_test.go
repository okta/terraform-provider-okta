package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaAuthServerClaimDefault_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthServerClaimDefault)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthServerClaimDefault, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAuthServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "address"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "SYSTEM"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "IDENTITY"),
					resource.TestCheckResourceAttr(resourceName, "always_include_in_token", "false"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "address"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "SYSTEM"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "IDENTITY"),
					resource.TestCheckResourceAttr(resourceName, "always_include_in_token", "true"),
				),
			},
		},
	})
}
