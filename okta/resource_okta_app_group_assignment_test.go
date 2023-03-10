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
	resourceName0 := fmt.Sprintf("%s.test.0", appGroupAssignment)
	resourceName1 := fmt.Sprintf("%s.test.1", appGroupAssignment)
	resourceName3 := fmt.Sprintf("%s.test3", appGroupAssignment)
	mgr := newFixtureManager(appGroupAssignment)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)

	// TF concurrency will cause this test to flap if the groups and assigned
	// priorities aren't executed in proper order
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "priority"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "priority"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName3),
					resource.TestCheckResourceAttrSet(resourceName3, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName3, "group_id"),
					resource.TestCheckResourceAttr(resourceName3, "priority", "3"),
					resource.TestCheckResourceAttr(resourceName3, "profile", "{}"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "priority"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "priority"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName3),
					resource.TestCheckResourceAttrSet(resourceName3, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName3, "group_id"),
					resource.TestCheckResourceAttr(resourceName3, "priority", "4"),
					resource.TestCheckResourceAttr(resourceName3, "profile", "{}"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "priority"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "priority"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName3),
					resource.TestCheckResourceAttrSet(resourceName3, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName3, "group_id"),
					resource.TestCheckResourceAttr(resourceName3, "priority", "3"),
					resource.TestCheckResourceAttr(resourceName3, "profile", "{}"),
				),
			},
			{
				ResourceName:      resourceName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName3]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName3)
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
	groupName := fmt.Sprintf("%s.test", group)
	mgr := newFixtureManager(appGroupAssignment)
	retainAssignment := mgr.GetFixtures("retain_assignment.tf", ri, t)
	retainAssignmentDestroy := mgr.GetFixtures("retain_assignment_destroy.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
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

func TestAccAppGroupAssignment_timeouts(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appGroupAssignment)
	resourceName0 := fmt.Sprintf("%s.test.0", appGroupAssignment)
	config := `
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"

  lifecycle {
    ignore_changes = [users, groups]
  }
}

resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

locals {
  group_ids = tolist([okta_group.test.id])
}

resource "okta_app_group_assignment" "test" {
  count = length(local.group_ids)

  app_id   = okta_app_oauth.test.id
  group_id = local.group_ids[count.index]
  priority = count.index

  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName0, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName0, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName0, "timeouts.update", "30m"),
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
