package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceIdpSaml_read(t *testing.T) {
	mgr := newFixtureManager(idpSaml, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	updatedConfig := mgr.GetFixtures("datasource_id.tf", t)
	idpSaml := mgr.GetFixtures("basic.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
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
