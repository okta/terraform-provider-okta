package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaIdpSocial_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", idpSocial, t.Name())
	preConfig := mgr.GetFixtures("basic.tf", t)
	config := mgr.GetFixtures("datasource.tf", t)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_idp_social.test_facebook", "id"),
					resource.TestCheckResourceAttrSet("data.okta_idp_social.test_google", "name"),
					resource.TestCheckResourceAttrSet("data.okta_idp_social.test_microsoft", "id"),
				),
			},
		},
	})
}
