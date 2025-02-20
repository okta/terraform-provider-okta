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

// Test creation of a simple AWS SWA app. The preconfigured apps are created by name.
func TestAccResourceOktaAppSwaApplication_preconfig(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSwa, t.Name())
	config := mgr.GetFixtures("preconfig.tf", t)
	updatedConfig := mgr.GetFixtures("preconfig_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSwa)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAppSwa, createDoesAppExist(sdk.NewSwaApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

// Test creation of a custom SAML app.
func TestAccResourceOktaAppSwaApplication_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSwa, t.Name())
	config := mgr.GetFixtures("custom.tf", t)
	updatedConfig := mgr.GetFixtures("custom_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSwa)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAppSwa, createDoesAppExist(sdk.NewSwaApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login.html"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login-updated.html"),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login-updated"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password-updated"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username-updated"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSwaApplication_timeouts(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSwa, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSwa)
	config := `
resource "okta_app_swa" "test" {
  label          = "testAcc_replace_with_uuid"
  button_field   = "btn-login"
  password_field = "txtbox-password"
  username_field = "txtbox-username"
  url            = "https://example.com/login.html"
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
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAppSwa, createDoesAppExist(sdk.NewSwaApplication())),
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
