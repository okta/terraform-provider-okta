package okta

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaAppUser_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", appUser)
	mgr := newFixtureManager(appUser, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	update := mgr.GetFixtures("update.tf", t)
	basicProfile := mgr.GetFixtures("basic_profile.tf", t)
	ignoreConfig := mgr.GetFixtures("profile_attributes_to_ignore.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkAppUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("testAcc_%d@example.com", mgr.Seed)),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					ensureAppUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", buildResourceName(mgr.Seed)),
				),
			},
			{
				Config: basicProfile,
				Check: resource.ComposeTestCheckFunc(
					ensureAppUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("testAcc_%d@example.com", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "profile", "{\"testCustom\":\"testing\"}"),
				),
			},
			{
				Config:  ignoreConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("testAcc_%d@example.com", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "profile", "{\"someAttribute\":\"testing\"}"), // Note: "testCustom" is ignored and should not be present
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName)
					}

					appID := rs.Primary.Attributes["app_id"]
					userID := rs.Primary.Attributes["user_id"]

					return fmt.Sprintf("%s/%s", appID, userID), nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import schema into state")
					}

					return nil
				},
			},
		},
	})
}

func TestAccOktaAppUser_retain(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", appUser)
	appName := fmt.Sprintf("%s.test", appOAuth)
	userName := fmt.Sprintf("%s.test", user)
	mgr := newFixtureManager(appUser, t.Name())
	retain := mgr.GetFixtures("retain.tf", t)
	retainDestroy := mgr.GetFixtures("retain_destroy.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkAppUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: retain,
				Check: resource.ComposeTestCheckFunc(
					ensureAppUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("testAcc_%d@example.com", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "retain_assignment", "true"),
				),
			},
			{
				Config: retainDestroy,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceNotExists(resourceName),
					ensureAppUserRetained(appName, userName),
				),
			},
		},
	})
}

func ensureAppUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", resourceName)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return missingErr
		}

		appID := rs.Primary.Attributes["app_id"]
		userID := rs.Primary.Attributes["user_id"]
		client := oktaClientForTest()

		u, _, err := client.Application.GetApplicationUser(context.Background(), appID, userID, nil)
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

		appID := rs.Primary.Attributes["app_id"]
		userID := rs.Primary.Attributes["user_id"]

		client := oktaClientForTest()
		_, response, err := client.Application.GetApplicationUser(context.Background(), appID, userID, nil)
		exists, err := doesResourceExist(response, err)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("resource still exists, App Id: %s, User Id: %s", appID, userID)
		}
	}

	return nil
}

func ensureAppUserRetained(appName, userName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		notFound := "resource not found: %s"
		// app user has been removed from state, so use app and user to query okta
		appRes, ok := s.RootModule().Resources[appName]
		if !ok {
			return fmt.Errorf(notFound, appName)
		}

		userRes, ok := s.RootModule().Resources[userName]
		if !ok {
			return fmt.Errorf(notFound, userName)
		}

		appID := appRes.Primary.ID
		userID := userRes.Primary.ID
		client := oktaClientForTest()

		g, _, err := client.Application.GetApplicationUser(context.Background(), appID, userID, nil)
		if err != nil {
			return err
		} else if g == nil {
			return fmt.Errorf("Application User not found for app ID, user ID: %s, %s", appID, userID)
		}

		return nil
	}
}
