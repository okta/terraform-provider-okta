package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaAppOauth_read(t *testing.T) {
	mgr := newFixtureManager(appOAuth, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	appCreate := buildTestAppOauth(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "client_secret"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "grant_types.#"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "redirect_uris.#"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "type"),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test", "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test_label", "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test", "status", statusActive),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test_label", "status", statusActive),
				),
			},
		},
	})
}

func buildTestAppOauth(d int) string {
	return fmt.Sprintf(`
resource "okta_app_oauth" "test" {
  label                      = "testAcc_%d"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["http://d.com/"]
  response_types             = ["code"]
  client_basic_secret        = "something_from_somewhere"
  client_id                  = "something_from_somewhere"
  token_endpoint_auth_method = "client_secret_basic"
  consent_method             = "TRUSTED"
}`, d)
}
