package okta

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceOktaSmtpServer_crud(t *testing.T) {
	mgr := newFixtureManager("resources", emailSmtp, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", emailSmtp)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "host"),
					resource.TestCheckResourceAttr(resourceName, "port", "8086"),
					resource.TestCheckResourceAttr(resourceName, "username", "test_user"),
				),
			},
		},
	})
}
