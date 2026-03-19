package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaAppUser_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppUser)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppUser, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	update := mgr.GetFixtures("update.tf", t)
	basicProfile := mgr.GetFixtures("basic_profile.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppUserDestroy,
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
					resource.TestCheckResourceAttr(resourceName, "username", acctest.BuildResourceName(mgr.Seed)),
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

func TestAccResourceOktaAppUser_retain(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppUser)
	appName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)
	userName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppUser, t.Name())
	retain := mgr.GetFixtures("retain.tf", t)
	retainDestroy := mgr.GetFixtures("retain_destroy.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppUserDestroy,
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
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()

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
		if rs.Type != resources.OktaIDaaSAppUser {
			continue
		}

		appID := rs.Primary.Attributes["app_id"]
		userID := rs.Primary.Attributes["user_id"]

		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
		_, response, err := client.Application.GetApplicationUser(context.Background(), appID, userID, nil)
		exists, err := utils.DoesResourceExist(response, err)
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
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()

		g, _, err := client.Application.GetApplicationUser(context.Background(), appID, userID, nil)
		if err != nil {
			return err
		} else if g == nil {
			return fmt.Errorf("Application User not found for app ID, user ID: %s, %s", appID, userID)
		}

		return nil
	}
}
