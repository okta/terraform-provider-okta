package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAppGroupAssignments_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", appGroupAssignments)
	mgr := newFixtureManager(appGroupAssignments, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)

	group1 := fmt.Sprintf("%s.test1", group)
	group2 := fmt.Sprintf("%s.test2", group)
	group3 := fmt.Sprintf("%s.test3", group)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
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
		},
	})
}

func ensureAppGroupAssignmentsExist(name string, groupsExpected ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		appID := rs.Primary.Attributes["app_id"]
		client := getOktaClientFromMetadata(testAccProvider.Meta())

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
