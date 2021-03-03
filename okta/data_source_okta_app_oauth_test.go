package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAppOauth_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	appCreate := buildTestAppOauth(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: appCreate,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "grant_types.#"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "redirect_uris.#"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "type"),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test_label", "label", buildResourceName(ri)),
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
