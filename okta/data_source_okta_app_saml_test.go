package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAppSaml_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_app_saml")
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_saml.test", "key_id"),
					/*resource.TestCheckResourceAttr("data.okta_app_saml.test", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app_saml.test", "status", statusActive),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "status", statusActive),
					*/),
			},
		},
	})
}
