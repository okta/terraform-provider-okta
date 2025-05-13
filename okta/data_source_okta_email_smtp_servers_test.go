package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceOktaEmailSmtpServers_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", emailSMTPServer, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_email_smtp_server.test", "id"),
					resource.TestCheckResourceAttr("data.okta_email_smtp_server.test", "host", "smtp.example.com"),
					resource.TestCheckResourceAttr("data.okta_email_smtp_server.test", "alias", "server4"),
					resource.TestCheckResourceAttr("data.okta_email_smtp_server.test", "username", "abcd"),
					resource.TestCheckResourceAttr("data.okta_email_smtp_server.test", "port", "587"),
				),
			},
		},
	})
}
