package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccDataSourceOktaAppOauth_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppOAuth, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	appCreate := buildTestAppOauth(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
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
					resource.TestCheckResourceAttr("data.okta_app_oauth.test", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test_label", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test", "status", idaas.StatusActive),
					resource.TestCheckResourceAttr("data.okta_app_oauth.test_label", "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "authentication_policy"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test_label", "authentication_policy"),
				),
			},
		},
	})
}

func TestAccDataSourceOktaAppOauth_authentication_policy(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppOAuth, t.Name())
	config := testOauthAuthenticationPolicyConfig(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_app_oauth.test", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "authentication_policy"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.okta_app_oauth.test", "type"),
				),
			},
		},
	})
}

func testOauthAuthenticationPolicyConfig(i int) string {
	return fmt.Sprintf(`
resource "okta_app_oauth" "test" {
  label          = "testAcc_%d"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
  consent_method = "TRUSTED"
}

data "okta_app_oauth" "test" {
  label = "testAcc_%d"
  depends_on = [okta_app_oauth.test]
}
`, i, i)
}

func buildTestAppOauth(d int) string {
	return fmt.Sprintf(`
resource "okta_app_oauth" "test" {
  label                      = "testAcc_%d"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["http://d.com/"]
  response_types             = ["code"]
  token_endpoint_auth_method = "client_secret_basic"
  consent_method             = "TRUSTED"
}`, d)
}
