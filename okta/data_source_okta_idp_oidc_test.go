package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceIdpOidc_read(t *testing.T) {
	mgr := newFixtureManager(idpOidc, t.Name())
	idpOidcConfig := mgr.GetFixtures("generic_oidc.tf", t)
	config := mgr.GetFixtures("datasource.tf", t)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
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
