package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaAppSecurePasswordStoreApplication_credsSchemes(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSecurePasswordStore, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSecurePasswordStore)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAppSecurePasswordStore, createDoesAppExist(sdk.NewSecurePasswordStoreApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSecurePasswordStoreApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "url", "http://test.com"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "pass"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "ADMIN_SETS_CREDENTIALS"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSecurePasswordStoreApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "url", "http://test.com"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "pass"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "EXTERNAL_PASSWORD_SYNC"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSecurePasswordStoreApplication_timeouts(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSecurePasswordStore, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSecurePasswordStore)
	config := `
resource "okta_app_secure_password_store" "test" {
  label              = "testAcc_replace_with_uuid"
  username_field     = "user"
  password_field     = "pass"
  url                = "http://test.com"
  credentials_scheme = "ADMIN_SETS_CREDENTIALS"
  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAppSecurePasswordStore, createDoesAppExist(sdk.NewSecurePasswordStoreApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}
