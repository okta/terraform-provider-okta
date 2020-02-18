package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

// Test creation of a simple AWS SWA app. The preconfigured apps are created by name.
func TestAccAppSwaApplication_preconfig(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSwa)
	config := mgr.GetFixtures("preconfig.tf", ri, t)
	updatedConfig := mgr.GetFixtures("preconfig_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appSwa)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSwa, createDoesAppExist(okta.NewSwaApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
		},
	})
}

// Test creation of a custom SAML app.
func TestAccAppSwaApplication_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSwa)
	config := mgr.GetFixtures("custom.tf", ri, t)
	updatedConfig := mgr.GetFixtures("custom_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appSwa)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSwa, createDoesAppExist(okta.NewSwaApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
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
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login-updated.html"),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login-updated"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password-updated"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username-updated"),
				),
			},
		},
	})
}
