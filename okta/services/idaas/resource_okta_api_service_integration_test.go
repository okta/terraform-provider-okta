package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAPIServiceIntegration_crd(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAPIServiceIntegration, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.example", resources.OktaIDaaSAPIServiceIntegration)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_api_service_integration.example", "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "anzennaapiservice"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "granted_scopes.*", map[string]string{
						"scope": "okta.users.read",
					}),
				),
			},
		},
	})
}
