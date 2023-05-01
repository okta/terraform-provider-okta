package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAuthenticator_read(t *testing.T) {
	mgr := newFixtureManager(authenticator, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := fmt.Sprintf("data.%s.test", authenticator)    // security question
	resourceName1 := fmt.Sprintf("data.%s.test_1", authenticator) // okta verify

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "settings"),
					resource.TestCheckNoResourceAttr(resourceName, "provider"),
					resource.TestCheckNoResourceAttr(resourceName, "provider_type"),
					resource.TestCheckNoResourceAttr(resourceName, "provider_hostname"),
					resource.TestCheckNoResourceAttr(resourceName, "provider_auth_port"),
					resource.TestCheckNoResourceAttr(resourceName, "provider_instance_id"),
					resource.TestCheckNoResourceAttr(resourceName, "provider_host"),
					resource.TestCheckNoResourceAttr(resourceName, "provider_secret_key"),
					resource.TestCheckNoResourceAttr(resourceName, "provider_integration_key"),

					resource.TestCheckResourceAttr(resourceName1, "type", "app"),
					resource.TestCheckResourceAttr(resourceName1, "key", "okta_verify"),
					resource.TestCheckResourceAttr(resourceName1, "name", "Okta Verify"),
					resource.TestCheckResourceAttrSet(resourceName1, "id"),
					resource.TestCheckResourceAttrSet(resourceName1, "status"),
					resource.TestCheckResourceAttrSet(resourceName1, "settings"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider_type"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider_hostname"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider_auth_port"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider_instance_id"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider_host"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider_secret_key"),
					resource.TestCheckNoResourceAttr(resourceName1, "provider_integration_key"),
				),
			},
		},
	})
}
