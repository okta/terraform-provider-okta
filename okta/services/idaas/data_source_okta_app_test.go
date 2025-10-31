package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccDataSourceOktaApp_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSApp, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	appCreate := buildTestApp(mgr.Seed)

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
					resource.TestCheckResourceAttrSet("okta_app_oauth.test", "id"),
					resource.TestCheckResourceAttr("data.okta_app.test", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app.test2", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app.test3", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app.test", "status", idaas.StatusActive),
					resource.TestCheckResourceAttr("data.okta_app.test2", "status", idaas.StatusActive),
					resource.TestCheckResourceAttr("data.okta_app.test3", "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet("data.okta_app.test", "authentication_policy"),
					resource.TestCheckResourceAttrSet("data.okta_app.test2", "authentication_policy"),
					resource.TestCheckResourceAttrSet("data.okta_app.test3", "authentication_policy"),
				),
			},
		},
	})
}

func buildTestApp(i int) string {
	return fmt.Sprintf(`
resource "okta_app_oauth" "test" {
  label          = "testAcc_%d"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
  consent_method = "TRUSTED"
}`, i)
}

func TestAccDataSourceOktaApp_label_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSApp, t.Name())
	config := testLabelConfig(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_app.test", "label", acctest.BuildResourceName(mgr.Seed)),
				),
			},
		},
	})
}

func TestAccDataSourceOktaApp_authentication_policy(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSApp, t.Name())
	config := testAuthenticationPolicyConfig(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_app.test", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttrSet("data.okta_app.test", "authentication_policy"),
				),
			},
		},
	})
}

func testAuthenticationPolicyConfig(i int) string {
	return fmt.Sprintf(`
resource "okta_app_oauth" "test" {
  label          = "testAcc_%d"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
  consent_method = "TRUSTED"
}

data "okta_app" "test" {
  label = "testAcc_%d"
  depends_on = [okta_app_oauth.test]
}
`, i, i)
}

func testLabelConfig(i int) string {
	return fmt.Sprintf(`
resource "okta_app_oauth" "test-dev" {
  label          = "testAcc_%d-dev"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
  consent_method = "TRUSTED"
}
resource "okta_app_oauth" "test" {
  label          = "testAcc_%d"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
  consent_method = "TRUSTED"
}
# before bug fix for #1111 data.okta_app.test wouldn't find okta_app_oauth.test
# correctly when it a had sibling with the same lable but additional information
# such as "myapp", and "myapp-dev"
data "okta_app" "test" {
  label = "testAcc_%d"
  depends_on = [
    okta_app_oauth.test-dev,
    okta_app_oauth.test
  ]
}
`, i, i, i)
}
