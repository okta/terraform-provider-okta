package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaAppGroupAssignments_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", appGroupAssignments)
	mgr := newFixtureManager("resources", appGroupAssignments, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)

	group1 := fmt.Sprintf("%s.test1", group)
	group2 := fmt.Sprintf("%s.test2", group)
	group3 := fmt.Sprintf("%s.test3", group)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentsExist(resourceName, group1, group2, group3),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentsExist(resourceName, group1, group2),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentsExist(resourceName, group1, group2, group3),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentsExist(resourceName, group1, group2, group3),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
		},
	})
}

// TestAccResourceOktaAppGroupAssignments_1088_unplanned_changes This test
// demonstrates incorrect behavior in okta_app_group_assignments has been
// corrected.  The original author implemented incorrect behavior, in terms of
// idiomatic design principles of TF providers, where it would proactively
// remove group assignments from an app if they were made outside of the
// resource. The correct behavior is to surface drift detection if a group is
// assigned to an app outside of this resource.
func TestAccResourceOktaAppGroupAssignments_1088_unplanned_changes(t *testing.T) {
	mgr := newFixtureManager("resources", appGroupAssignments, t.Name())
	assignments1 := fmt.Sprintf("%s.test", appGroupAssignments)
	bookmarkApp := fmt.Sprintf("%s.test", appBookmark)
	groupA := fmt.Sprintf("%s.a", group)
	groupB := fmt.Sprintf("%s.b", group)
	groupC := fmt.Sprintf("%s.c", group)

	baseConfig := `
resource "okta_app_bookmark" "test" {
	label = "testAcc_replace_with_uuid"
	url   = "https://test.com"
}
resource "okta_group" "a" {
	name        = "testAcc-group-a_replace_with_uuid"
	description = "Group A"
}
resource "okta_group" "b" {
	name        = "testAcc-group-b_replace_with_uuid"
	description = "Group B"
}
resource "okta_group" "c" {
	name        = "testAcc-group-c_replace_with_uuid"
	description = "Group C"
}`

	step1Config := `
resource "okta_app_group_assignments" "test" {
	app_id = okta_app_bookmark.test.id
	group {
		id = okta_group.a.id
		priority = 1
		profile = jsonencode({"test": "a"})
	}
}`

	step2Config := `
resource "okta_app_group_assignments" "test" {
	app_id = okta_app_bookmark.test.id
	group {
		id = okta_group.a.id
		priority = 1
		profile = jsonencode({"test": "a"})
	}
	group {
		id = okta_group.b.id
		priority = 2
		profile = jsonencode({"test": "b"})
	}
}`

	step3Config := `
resource "okta_app_group_assignments" "test" {
	app_id = okta_app_bookmark.test.id
	group {
		id = okta_group.a.id
		priority = 1
		profile = jsonencode({"test": "a"})
	}
	group {
		id = okta_group.b.id
		priority = 2
		profile = jsonencode({"test": "b"})
	}
	group {
		id = okta_group.c.id
		priority = 4
		profile = jsonencode({"test": "c"})
	}
}`

	stepLastConfig := `
resource "okta_app_group_assignments" "test" {
	app_id = okta_app_bookmark.test.id
	group {
		id = okta_group.a.id
		priority = 99
		profile = jsonencode({"different": "profile value"})
	}
}`

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		CheckDestroy:      nil,
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Vanilla step
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step1Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(assignments1, "group.#", "1"),
					resource.TestCheckResourceAttr(assignments1, "group.0.priority", "1"),
					resource.TestCheckResourceAttr(assignments1, "group.0.profile", `{"test":"a"}`),
					ensureAppGroupAssignmentsExist(assignments1, groupA),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step2Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(assignments1, "group.#", "2"),
					resource.TestCheckResourceAttr(assignments1, "group.0.priority", "1"),
					resource.TestCheckResourceAttr(assignments1, "group.0.profile", `{"test":"a"}`),
					resource.TestCheckResourceAttr(assignments1, "group.1.priority", "2"),
					resource.TestCheckResourceAttr(assignments1, "group.1.profile", `{"test":"b"}`),
					ensureAppGroupAssignmentsExist(assignments1, groupA, groupB),

					// This mimics assigning Group C to the app outside of
					// Terraform. In this case doing so with a direct API call
					// via the test harness which is equivalent to "Click Ops"
					clickOpsAssignGroupToApp(bookmarkApp, groupC),
					clickOpsCheckIfGroupIsAssignedToApp(bookmarkApp, groupA, groupB, groupC),

					// NOTE: after these checks run the terraform test runner
					// will do a refresh and catch that group C has been added
					// to the app outside of the terraform config and emit a
					// non-empty plan
				),

				// side effect of the TF test runner is expecting a non-empty
				// plan is treated as an apply accept and adds group c to local
				// state
				ExpectNonEmptyPlan: true,
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step3Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(assignments1, "group.#", "3"),
					resource.TestCheckResourceAttr(assignments1, "group.0.priority", "1"),
					resource.TestCheckResourceAttr(assignments1, "group.0.profile", `{"test":"a"}`),
					resource.TestCheckResourceAttr(assignments1, "group.1.priority", "2"),
					resource.TestCheckResourceAttr(assignments1, "group.1.profile", `{"test":"b"}`),
					resource.TestCheckResourceAttr(assignments1, "group.2.priority", "4"),
					resource.TestCheckResourceAttr(assignments1, "group.2.profile", `{"test":"c"}`),
					ensureAppGroupAssignmentsExist(assignments1, groupA, groupB, groupC),
				),
			},
			{
				// check that we can do removing group assignments
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step2Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(assignments1, "group.#", "2"),
					resource.TestCheckResourceAttr(assignments1, "group.0.priority", "1"),
					resource.TestCheckResourceAttr(assignments1, "group.0.profile", `{"test":"a"}`),
					resource.TestCheckResourceAttr(assignments1, "group.1.priority", "2"),
					resource.TestCheckResourceAttr(assignments1, "group.1.profile", `{"test":"b"}`),
					ensureAppGroupAssignmentsExist(assignments1, groupA, groupB),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step1Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(assignments1, "group.#", "1"),
					resource.TestCheckResourceAttr(assignments1, "group.0.priority", "1"),
					resource.TestCheckResourceAttr(assignments1, "group.0.profile", `{"test":"a"}`),
					ensureAppGroupAssignmentsExist(assignments1, groupA),
				),
			},
			{
				// Check that priority and profile can be changed on the group
				// itself
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, stepLastConfig)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(assignments1, "group.#", "1"),
					resource.TestCheckResourceAttr(assignments1, "group.0.priority", "99"),
					resource.TestCheckResourceAttr(assignments1, "group.0.profile", `{"different":"profile value"}`),
					ensureAppGroupAssignmentsExist(assignments1, groupA),
				),
			},
		},
	})
}

func ensureAppGroupAssignmentsExist(resourceName string, groupsExpected ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", resourceName)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return missingErr
		}

		appID := rs.Primary.Attributes["app_id"]
		client := sdkV2ClientForTest()

		// Get all the IDs of groups we expect to be assigned
		expectedGroupIDs := map[string]bool{}
		for _, groupExpected := range groupsExpected {
			grs, ok := s.RootModule().Resources[groupExpected]
			if !ok {
				return missingErr
			}
			expectedGroupIDs[grs.Primary.Attributes["id"]] = false
		}

		for i := 0; i < len(groupsExpected); i++ {
			groupID := rs.Primary.Attributes[fmt.Sprintf("group.%d.id", i)]
			g, _, err := client.Application.GetApplicationGroupAssignment(context.Background(), appID, groupID, nil)
			if err != nil {
				return err
			} else if g == nil {
				return missingErr
			}
			// group found, check it off
			expectedGroupIDs[groupID] = true
		}

		// now check we found all the groupIDs we expected
		if len(expectedGroupIDs) != len(groupsExpected) {
			return fmt.Errorf("expected %d assigned groups but got %d", len(groupsExpected), len(expectedGroupIDs))
		}

		// make sure we found them all
		for groupID, found := range expectedGroupIDs {
			if !found {
				return fmt.Errorf("expected group %s to be assigned but wasn't", groupID)
			}
		}
		return nil
	}
}

func clickOpsAssignGroupToApp(appResourceName, groupResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resources := []string{appResourceName, groupResourceName}
		for _, resourceName := range resources {
			missingErr := fmt.Errorf("resource not found: %s", resourceName)
			if _, ok := s.RootModule().Resources[resourceName]; !ok {
				return missingErr
			}
		}

		appRS := s.RootModule().Resources[appResourceName]
		appID := appRS.Primary.Attributes["id"]
		groupRS := s.RootModule().Resources[groupResourceName]
		groupID := groupRS.Primary.Attributes["id"]
		client := sdkV2ClientForTest()
		_, _, err := client.Application.CreateApplicationGroupAssignment(context.Background(), appID, groupID, sdk.ApplicationGroupAssignment{})
		if err != nil {
			return fmt.Errorf("API: unable to assign app %q to group %q, err: %+v", appID, groupID, err)
		}

		return nil
	}
}

func clickOpsCheckIfGroupIsAssignedToApp(appResourceName string, groups ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, groupResourceName := range groups {
			resources := []string{appResourceName, groupResourceName}
			for _, resourceName := range resources {
				missingErr := fmt.Errorf("resource not found: %s", resourceName)
				if _, ok := s.RootModule().Resources[resourceName]; !ok {
					return missingErr
				}
			}

			appRS := s.RootModule().Resources[appResourceName]
			appID := appRS.Primary.Attributes["id"]
			groupRS := s.RootModule().Resources[groupResourceName]
			groupID := groupRS.Primary.Attributes["id"]
			client := sdkV2ClientForTest()
			_, _, err := client.Application.GetApplicationGroupAssignment(context.Background(), appID, groupID, nil)
			if err != nil {
				return fmt.Errorf("API: app %q is not assigned to group %s", appID, groupID)
			}
		}

		return nil
	}
}

// This test demonstrate the ability to unassigned all groups from app without having to destroy the resource
// This behavior is already enabled by the API
func TestAccResourceOktaAppGroupAssignments_2068_empty_assignments(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", appGroupAssignments)
	mgr := newFixtureManager("resources", appGroupAssignments, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated_empty.tf", t)

	group1 := fmt.Sprintf("%s.test1", group)
	group2 := fmt.Sprintf("%s.test2", group)
	group3 := fmt.Sprintf("%s.test3", group)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentsExist(resourceName, group1, group2, group3),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttr(resourceName, "group.#", "3"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttr(resourceName, "group.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppBGroupAssignments_timeouts(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", appGroupAssignments)
	mgr := newFixtureManager("resources", appGroupAssignments, t.Name())
	bookmarkApp := fmt.Sprintf("%s.test", appBookmark)
	groupA := fmt.Sprintf("%s.a", group)
	groupB := fmt.Sprintf("%s.b", group)
	config := `
resource "okta_app_bookmark" "test" {
	label = "testAcc_replace_with_uuid"
	url   = "https://test.com"
}
resource "okta_group" "a" {
	name        = "testAcc-group-a_replace_with_uuid"
	description = "Group A"
}
resource "okta_group" "b" {
	name        = "testAcc-group-b_replace_with_uuid"
	description = "Group B"
}
resource "okta_app_group_assignments" "test" {
  app_id = okta_app_bookmark.test.id
  group {
    id = okta_group.a.id
  }
  group {
    id = okta_group.b.id
  }
  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appGroupAssignments, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),

					clickOpsCheckIfGroupIsAssignedToApp(bookmarkApp, groupA, groupB),
				),
			},
		},
	})
}
