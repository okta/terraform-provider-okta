package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// Test creation of a simple AWS WSFederation app. The pre-configured apps are created by name.
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
					resource.TestCheckResourceAttrSet(resourceName, "site_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewWsFederationApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttrSet(resourceName, "site_url"),
				),
			},
		},
	})
}

// Test creation of a custom WSFed app.
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
					resource.TestCheckResourceAttr(resourceName, "realm", "test"),
					resource.TestCheckResourceAttr(resourceName, "name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"),
					resource.TestCheckResourceAttr(resourceName, "audience_restriction", "https://signin.test.com"),
					resource.TestCheckResourceAttr(resourceName, "authn_context_class_ref", "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"),
					resource.TestCheckResourceAttr(resourceName, "group_filter", "app1.*"),
					resource.TestCheckResourceAttr(resourceName, "group_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "group_value_format", "dn"),
					resource.TestCheckResourceAttr(resourceName, "username_attribute", "username"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|bob|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|hope|"),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewWsFederationApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "realm", "test"),
					resource.TestCheckResourceAttr(resourceName, "name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"),
					resource.TestCheckResourceAttr(resourceName, "audience_restriction", "https://signin.test.com"),
					resource.TestCheckResourceAttr(resourceName, "authn_context_class_ref", "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"),
					resource.TestCheckResourceAttr(resourceName, "group_filter", "app1.*"),
					resource.TestCheckResourceAttr(resourceName, "group_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "group_value_format", "dn"),
					resource.TestCheckResourceAttr(resourceName, "username_attribute", "username"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|bob|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|hope|"),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
				),
			},
		},
	})
}

func TestAccAppWsFedApplication_timeouts(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appWsFed)
	resourceName := fmt.Sprintf("%s.test", appWsFed)
	importConfig := mgr.GetFixtures("import.tf", ri, t)
	config := `
	 resource "okta_app_ws_federation" "test" {
		    label    = "testAcc_replace_with_uuid"
		    site_url = "https://signin.example.com/saml"
		    reply_url = "https://example.com"
		    reply_override = false
			realm = "test"
		    name_id_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
		    audience_restriction = "https://signin.example.com"
		    authn_context_class_ref = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
		    group_filter = "app1.*"
		    group_name = "username"
		    group_value_format = "dn"
		    username_attribute = "username"
		    attribute_statements = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|bob|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|hope|"
		    visibility = false
		    status = "ACTIVE"

			timeouts {
				create = "60m"
				read = "2h"
				update = "30m"
			}			
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
			{
				Config: importConfig,
			},
		},
	})
}
