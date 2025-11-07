package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAPIServiceIntegration_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAPIServiceIntegration, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_api_service_integration.example", "type", "anzennaapiservice"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.okta_api_service_integration.example",
						"granted_scopes.*",
						map[string]string{
							"scope": "okta.users.read",
						},
					),
				),
			},
		},
	})
}
