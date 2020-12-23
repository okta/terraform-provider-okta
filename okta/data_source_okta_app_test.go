package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceApp_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_app")
	config := mgr.GetFixtures("datasource.tf", ri, t)
	appCreate := buildTestApp(ri)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: appCreate,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_app_oauth.test", "id"),
					resource.TestCheckResourceAttr("data.okta_app.test", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app.test2", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app.test3", "label", fmt.Sprintf("testAcc_%d", ri)),
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
