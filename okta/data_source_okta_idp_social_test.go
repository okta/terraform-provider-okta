package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceIdpSocial_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(idpSocial)
	preConfig := mgr.GetFixtures("basic.tf", ri, t)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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
