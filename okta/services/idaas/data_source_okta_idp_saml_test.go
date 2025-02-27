package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaIdpSaml_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSIdpSaml, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	updatedConfig := mgr.GetFixtures("datasource_id.tf", t)
	idpSaml := mgr.GetFixtures("basic.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: idpSaml,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_idp_saml.test", "id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_idp_saml.test", "id"),
				),
			},
		},
	})
}
