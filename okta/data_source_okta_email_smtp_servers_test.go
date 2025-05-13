package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceOktaEmailSmtpServers_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", emailSmtp, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_email_smtp.test", "id"),
					resource.TestCheckResourceAttr("data.okta_email_smtp.test", "host", "smtp.example.com"),
					resource.TestCheckResourceAttr("data.okta_email_smtp.test", "alias", "server4"),
					resource.TestCheckResourceAttr("data.okta_email_smtp.test", "username", "abcd"),
					resource.TestCheckResourceAttr("data.okta_email_smtp.test", "port", "587"),
				),
			},
		},
	})
}
