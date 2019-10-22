package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourceApp_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_app")
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_app_oauth.test", "id"),
					resource.TestCheckResourceAttr("data.okta_app.test", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app.test2", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app.test3", "label", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr("data.okta_app.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.okta_app.test2", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.okta_app.test3", "status", "ACTIVE"),
				),
			},
		},
	})
}
