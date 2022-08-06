package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// Test creation of a simple AWS SWA app. The preconfigured apps are created by name.
func TestAccAppSwaApplication_preconfig(t *testing.T) {
	mgr := newFixtureManager(appSwa, t.Name())
	config := mgr.GetFixtures("preconfig.tf", t)
	updatedConfig := mgr.GetFixtures("preconfig_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", appSwa)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appSwa, createDoesAppExist(okta.NewSwaApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

// Test creation of a custom SAML app.
func TestAccAppSwaApplication_crud(t *testing.T) {
	mgr := newFixtureManager(appSwa, t.Name())
	config := mgr.GetFixtures("custom.tf", t)
	updatedConfig := mgr.GetFixtures("custom_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", appSwa)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appSwa, createDoesAppExist(okta.NewSwaApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login.html"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login-updated.html"),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login-updated"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password-updated"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username-updated"),
				),
			},
		},
	})
}
