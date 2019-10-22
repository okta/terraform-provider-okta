package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourceIdpSaml_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_idp_saml")
	config := mgr.GetFixtures("datasource.tf", ri, t)
	updatedConfig := mgr.GetFixtures("datasource_id.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
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
