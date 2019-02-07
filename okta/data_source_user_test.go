package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaUser)
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
					resource.TestCheckResourceAttrSet("data.okta_user.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_user.test2", "id"),
					resource.TestCheckResourceAttr("data.okta_user.test", "first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_user.test", "last_name", "Smith"),
					resource.TestCheckResourceAttr("data.okta_user.test2", "first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_user.test2", "last_name", "Smith"),
					resource.TestCheckResourceAttr("data.okta_user.test", "status", "PROVISIONED"),
					resource.TestCheckResourceAttr("data.okta_user.test2", "status", "PROVISIONED"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("okta_user.testAcc_%d", ri), "id"),
				),
			},
		},
	})
}
