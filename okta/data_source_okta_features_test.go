package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODU
func TestAccDataSourceOktaFeature_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", app, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	appCreate := buildTestApp(mgr.Seed)

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
					resource.TestCheckResourceAttrSet("okta_app_oauth.test", "id"),
					resource.TestCheckResourceAttr("data.okta_app.test", "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app.test2", "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app.test3", "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app.test", "status", statusActive),
					resource.TestCheckResourceAttr("data.okta_app.test2", "status", statusActive),
					resource.TestCheckResourceAttr("data.okta_app.test3", "status", statusActive),
				),
			},
		},
	})
}
