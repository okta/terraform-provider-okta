package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAppWsFed_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appWsFed)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	appCreate := buildTestAppWsFed(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: appCreate,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "site_url", "https://signin.test.com/saml"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "reply_url", "https://test.com"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "realm", "test"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "audience_restriction", "https://signin.test.com"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "authn_context_class_ref", "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "group_filter", "app1.*"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "group_name", "username"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "group_value_format", "dn"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "username_attribute", "username"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "attribute_statements", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|bob|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|hope|"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "status", "ACTIVE"),
				),
			},
		},
	})
}

func buildTestAppWsFed(d int) string {
	return fmt.Sprintf(`
	resource "okta_app_ws_federation" "test" {
		label    = "testAcc_%d"
		site_url = "https://signin.test.com/saml"
		reply_url = "https://test.com"
		reply_override = false
		realm = "test"
		name_id_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
		audience_restriction = "https://signin.test.com"
		authn_context_class_ref = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
		group_filter = "app1.*"
		group_name = "username"
		group_value_format = "dn"
		username_attribute = "username"
		attribute_statements = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|bob|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|hope|"
		visibility = false
		status = "ACTIVE"
	}`, d)
}
