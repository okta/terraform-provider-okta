package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourceAppSaml_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_app_saml")
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_app_saml.test", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app_saml.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "status", "ACTIVE"),
				),
			},
		},
	})
}
