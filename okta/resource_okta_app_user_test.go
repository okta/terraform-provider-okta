package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOktaAppUser_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", appUser)
	mgr := newFixtureManager(appUser)
	config := mgr.GetFixtures("basic.tf", ri, t)
	update := mgr.GetFixtures("update.tf", ri, t)
	basicProfile := mgr.GetFixtures("basic_profile.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkAppUserDestroy,
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
			{
				Config: basicProfile,
				Check: resource.ComposeTestCheckFunc(
					ensureAppUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("testAcc_%d@example.com", ri)),
					resource.TestCheckResourceAttr(resourceName, "profile", "{\"testCustom\":\"testing\"}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName)
					}

					appId := rs.Primary.Attributes["app_id"]
					userId := rs.Primary.Attributes["user_id"]

					return fmt.Sprintf("%s/%s", appId, userId), nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("Failed to import schema into state")
					}

					return nil
				},
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

func checkAppUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != appUser {
			continue
		}

		appId := rs.Primary.Attributes["app_id"]
		userId := rs.Primary.Attributes["user_id"]

		client := getOktaClientFromMetadata(testAccProvider.Meta())
		_, response, err := client.Application.GetApplicationUser(appId, userId, nil)
		exists, err := doesResourceExist(response, err)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("resource still exists, App Id: %s, User Id: %s", appId, userId)
		}
	}

	return nil
}
