package okta

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceUsers(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_users")
	config := mgr.GetFixtures("basic.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Ensure users are created
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test2", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test3", "id"),
				),
			},
			{
				// Ensure data source props are set
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "3"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.1.first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.2.first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.last_name", "Doe"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.1.last_name", "Jones"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.2.last_name", "Entwhistle"),
				),
			},
		},
	})
}
