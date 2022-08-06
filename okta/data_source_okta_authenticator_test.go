package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAuthenticator_read(t *testing.T) {
	mgr := newFixtureManager(authenticator, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := fmt.Sprintf("data.%s.test", authenticator)
	resourceName1 := fmt.Sprintf("data.%s.test_1", authenticator)

	oktaResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "key"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "settings"),
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
					resource.TestCheckResourceAttrSet(resourceName1, "id"),
					resource.TestCheckResourceAttrSet(resourceName1, "key"),
					resource.TestCheckResourceAttrSet(resourceName1, "name"),
					resource.TestCheckResourceAttrSet(resourceName1, "status"),
					resource.TestCheckResourceAttrSet(resourceName1, "settings"),
					resource.TestCheckResourceAttr(resourceName1, "type", "app"),
					resource.TestCheckResourceAttr(resourceName1, "key", "okta_verify"),
					resource.TestCheckResourceAttr(resourceName1, "name", "Okta Verify"),
				),
			},
		},
	})
}
