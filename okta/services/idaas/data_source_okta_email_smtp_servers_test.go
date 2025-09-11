package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEmailSMTPServers_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSEmailSMTPServer, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
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
