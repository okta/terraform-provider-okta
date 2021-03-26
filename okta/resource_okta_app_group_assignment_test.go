package okta

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAppGroupAssignment_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", appGroupAssignment)
	resourceName0 := fmt.Sprintf("%s.test.0", appGroupAssignment)
	resourceName1 := fmt.Sprintf("%s.test.1", appGroupAssignment)
	mgr := newFixtureManager(appGroupAssignment)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttr(resourceName0, "priority", "0"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttr(resourceName1, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttr(resourceName0, "priority", "0"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttr(resourceName1, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttr(resourceName0, "priority", "0"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttr(resourceName1, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName0]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName)
					}
					appID := rs.Primary.Attributes["app_id"]
					groupID := rs.Primary.Attributes["group_id"]
					return fmt.Sprintf("%s/%s", appID, groupID), nil
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

func TestAccAppGroupAssignment_retain(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", appGroupAssignment)
	appName := fmt.Sprintf("%s.test", appOAuth)
	groupName := fmt.Sprintf("%s.test", oktaGroup)
	mgr := newFixtureManager(appGroupAssignment)
	retainAssignment := mgr.GetFixtures("retain_assignment.tf", ri, t)
	retainAssignmentDestroy := mgr.GetFixtures("retain_assignment_destroy.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: retainAssignment,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttr(resourceName, "retain_assignment", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile", "{}"),
				),
			},
			{
				Config: retainAssignmentDestroy,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceNotExists(resourceName),
					ensureAppGroupAssignmentRetained(appName, groupName),
				),
			},
		},
	})
}

func ensureAppGroupAssignmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		appID := rs.Primary.Attributes["app_id"]
		groupID := rs.Primary.Attributes["group_id"]
		client := getOktaClientFromMetadata(testAccProvider.Meta())

		g, _, err := client.Application.GetApplicationGroupAssignment(context.Background(), appID, groupID, nil)
		if err != nil {
			return err
		} else if g == nil {
			return missingErr
		}

		return nil
	}
}

func ensureAppGroupAssignmentRetained(appName, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		notFound := "resource not found: %s"
		// app group assignment has been removed from state, so use app and group to query okta
		appRes, ok := s.RootModule().Resources[appName]
		if !ok {
			return fmt.Errorf(notFound, appName)
		}

		groupRes, ok := s.RootModule().Resources[groupName]
		if !ok {
			return fmt.Errorf(notFound, groupName)
		}

		appID := appRes.Primary.ID
		groupID := groupRes.Primary.ID
		client := getOktaClientFromMetadata(testAccProvider.Meta())

		g, _, err := client.Application.GetApplicationGroupAssignment(context.Background(), appID, groupID, nil)
		if err != nil {
			return err
		} else if g == nil {
			return fmt.Errorf("application group assignment not found for app ID, group ID: %s, %s", appID, groupID)
		}
		return nil
	}
}
