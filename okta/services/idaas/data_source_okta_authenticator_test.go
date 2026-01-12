package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAuthenticator_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthenticator, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSAuthenticator)    // security question
	resourceName1 := fmt.Sprintf("data.%s.test_1", resources.OktaIDaaSAuthenticator) // okta verify

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
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

// TestAccDataSourceOktaAuthenticator_WithMethods tests data source reading of authenticators with methods
func TestAccDataSourceOktaAuthenticator_WithMethods(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticator)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSAuthenticator)
	config := `
resource "okta_authenticator" "test" {
  name   = "Phone"
  key    = "phone_number"
  status = "ACTIVE"

  method {
    type   = "sms"
    status = "ACTIVE"
  }

  method {
    type   = "voice"
    status = "INACTIVE"
  }
}

data "okta_authenticator" "test" {
  key = okta_authenticator.test.key
  depends_on = [okta_authenticator.test]
}
`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Verify resource was created
					resource.TestCheckResourceAttr(resourceName, "key", "phone_number"),
					resource.TestCheckResourceAttr(resourceName, "name", "Phone"),
					resource.TestCheckResourceAttr(resourceName, "method.#", "2"),

					// Verify data source reads the methods correctly
					resource.TestCheckResourceAttr(dataSourceName, "key", "phone_number"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "Phone"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "phone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttr(dataSourceName, "method.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "method.*", map[string]string{
						"type":   "sms",
						"status": "ACTIVE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "method.*", map[string]string{
						"type":   "voice",
						"status": "INACTIVE",
					}),
				),
			},
		},
	})
}
