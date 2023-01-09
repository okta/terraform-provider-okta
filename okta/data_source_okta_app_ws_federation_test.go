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
					resource.TestCheckResourceAttrSet("data.okta_app_ws_federation.test", "key_id"),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test_label", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test", "status", statusActive),
					resource.TestCheckResourceAttr("data.okta_app_ws_federation.test_label", "status", statusActive),
				),
			},
		},
	})
}

func buildTestAppWsFed(d int) string {
	return fmt.Sprintf(`
	resource "okta_app_ws_federation" "example" {
		label    = "example_%d"
		site_url = "https://signin.example.com/saml"
		realm = "example"
		reply_url = "https://example.com"
		allow_override = false
		name_id_format = "uid"
		audience_restriction = "https://signin.example.com"
		authn_context_class_ref = "Kerberos"
		group_filter = "app1.*"
		group_name = "username"
		group_value_format = "dn"
		username_attribute = "username"
		attribute_statements = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|${user.firstName}|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|${user.lastName}|"
		visibility = true      
	}`, d)
}
