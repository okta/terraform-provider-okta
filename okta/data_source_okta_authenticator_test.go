package okta

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAuthenticator_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaAuthenticator)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	configInvalid := mgr.GetFixtures("datasource_not_found.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_authenticator.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_authenticator.test", "key"),
					resource.TestCheckResourceAttrSet("data.okta_authenticator.test", "name"),
					resource.TestCheckResourceAttrSet("data.okta_authenticator.test", "status"),
					resource.TestCheckResourceAttrSet("data.okta_authenticator.test", "settings"),
					resource.TestCheckResourceAttr("data.okta_authenticator.test", "type", "security_question"),
					resource.TestCheckResourceAttr("data.okta_authenticator.test", "key", "security_question"),
					resource.TestCheckResourceAttr("data.okta_authenticator.test", "name", "Security Question"),
				),
			},
			{
				Config:      configInvalid,
				ExpectError: regexp.MustCompile(`\bdoes not exist`),
			},
		},
	})
}
