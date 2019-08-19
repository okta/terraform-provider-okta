package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppUser_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", appUser)
	mgr := newFixtureManager(appUser)
	config := mgr.GetFixtures("basic.tf", ri, t)
	update := mgr.GetFixtures("update.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("testAcc_%d@example.com", ri)),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					ensureAppUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("testAcc_%d", ri)),
				),
			},
		},
	})
}

func ensureAppUserExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		appId := rs.Primary.Attributes["app_id"]
		userId := rs.Primary.Attributes["user_id"]
		client := getOktaClientFromMetadata(testAccProvider.Meta())

		u, _, err := client.Application.GetApplicationUser(appId, userId, nil)
		if err != nil {
			return err
		} else if u == nil {
			return missingErr
		}

		return nil
	}
}
