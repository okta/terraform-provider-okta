package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceApp_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(app)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	appCreate := buildTestApp(ri)

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
					resource.TestCheckResourceAttrSet("okta_app_oauth.test", "id"),
					resource.TestCheckResourceAttr("data.okta_app.test", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app.test2", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app.test3", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app.test", "status", statusActive),
					resource.TestCheckResourceAttr("data.okta_app.test2", "status", statusActive),
					resource.TestCheckResourceAttr("data.okta_app.test3", "status", statusActive),
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

func TestAccOktaDataSourceAppLabelTest_read(t *testing.T) {
	ri := acctest.RandInt()
	config := testLabelConfig(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_app.test", "label", buildResourceName(ri)),
				),
			},
		},
	})
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
