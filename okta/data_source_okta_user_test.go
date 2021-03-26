package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceUser_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaUser)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	createUser := mgr.GetFixtures("datasource_create_user.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
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
					resource.TestCheckResourceAttr("data.okta_user.read_by_id", "first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_user.read_by_id", "last_name", "Smith"),
				),
			},
		},
	})
}
