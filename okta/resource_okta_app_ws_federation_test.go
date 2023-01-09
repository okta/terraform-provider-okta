package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// Test creation of a simple AWS WSFederation app. The preconfigured apps are created by name.
func TestAccAppWsFedApplication_preconfig(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appWsFed)
	config := mgr.GetFixtures("preconfig.tf", ri, t)
	updatedConfig := mgr.GetFixtures("preconfig_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appWsFed)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appWsFed, createDoesAppExist(okta.NewWsFederationApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewWsFederationApplication())),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewWsFederationApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

// Test creation of a custom SAML app.
func TestAccAppWsFedApplication_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appWsFed)
	config := mgr.GetFixtures("custom.tf", ri, t)
	updatedConfig := mgr.GetFixtures("custom_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appWsFed)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appWsFed, createDoesAppExist(okta.NewWsFederationApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewWsFederationApplication())),
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
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewWsFederationApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
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

func TestAccAppWsFedApplication_timeouts(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appWsFed)
	resourceName := fmt.Sprintf("%s.test", appWsFed)
	config := `
	resource "okta_app_ws_federation" "example" {
		label    = "example"
		site_url = "https://signin.example.com/saml"
		realm = "example"
		reply_to_url = "https://example.com"
		allow_reply_to_override = false
		name_id_format = "uid"
		audience_restriction = "https://signin.example.com"
		assert_authentication_context = "Kerberos"
		group_filter = "app1.*"
		group_attribute_name = "username"
		group_attribute_value = "dn"
		username_attribute = "username"
		custom_attribute_statements = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|${user.firstName}|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|${user.lastName}|"
		visibility = true      
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appWsFed, createDoesAppExist(okta.NewWsFederationApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}
