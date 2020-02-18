package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourceUsers_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_users")
	users := mgr.GetFixtures("users.tf", ri, t)
	config := mgr.GetFixtures("basic.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Ensure users are created
				Config: users,
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
					resource.TestCheckResourceAttrSet("data.okta_users.test", "users.#"),
				),
			},
		},
	})
}
