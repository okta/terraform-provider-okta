package okta

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceOktaSmtpServer_crud(t *testing.T) {
	mgr := newFixtureManager("resources", emailSMTPServer, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", emailSMTPServer)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "host", "192.168.2.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "8086"),
					resource.TestCheckResourceAttr(resourceName, "username", "test_user"),
					resource.TestCheckResourceAttr(resourceName, "alias", "CustomisedServer"),
				),
			},
		},
	})
}
