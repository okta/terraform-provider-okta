package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourceUser_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaUser)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	// Avoiding race conditions
	createUser := mgr.GetFixtures("datasource_create_user.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_user.test", "id"),
					resource.TestCheckResourceAttr("data.okta_user.test", "first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_user.test", "last_name", "Smith"),
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
				),
			},
		},
	})
}
