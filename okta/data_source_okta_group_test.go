package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOktaDataSourceGroup_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaGroup)
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
					resource.TestCheckResourceAttrSet("data.okta_group.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_group.test", "type"),
					resource.TestCheckResourceAttrSet("okta_group.test", "id"),
					resource.TestCheckResourceAttr("okta_group.test", "users.#", "1"),
				),
			},
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(testDoesNotExist("okta_group.test_type")),
			},
		},
	})
}

func testDoesNotExist(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if ok {
			return fmt.Errorf("Resource should not exist: %s", name)
		}
		return nil
	}
}
