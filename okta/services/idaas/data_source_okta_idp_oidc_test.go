package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaIdpOidc_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSIdpOidc, t.Name())
	idpOidcConfig := mgr.GetFixtures("generic_oidc.tf", t)
	config := mgr.GetFixtures("datasource.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: idpOidcConfig,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_idp_oidc.test", "id"),
				),
			},
		},
	})
}
