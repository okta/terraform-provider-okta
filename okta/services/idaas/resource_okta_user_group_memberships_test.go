package idaas_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// TestAccResourceOktaUserGroupMemberships_GH1218 addresses https://github.com/okta/terraform-provider-okta/issues/1218
// It tests the new track_all_groups attribute:
//   - track_all_groups = true makes Terraform track all of the user's group memberships,
//     detecting drift from groups added or removed outside of Terraform.
//   - Import with "user_id/true" correctly sets track_all_groups and populates groups from the API.
func TestAccResourceOktaUserGroupMemberships_GH1218(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUserGroupMemberships, t.Name())
	config := mgr.GetFixtures("basic_gh1218.tf", t)
	configUpdate := mgr.GetFixtures("basic_gh1218_update.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUserGroupMemberships)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				// Create user assigned to 2 groups with track_all_groups = true.
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "track_all_groups", "true"),
				),
			},
			{
				// Verify idempotency: a second plan produces no diff (no phantom drift
				// from Okta's built-in "Everyone" group or other BUILT_IN groups).
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Add a third managed group; groups.# should become 3.
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "groups.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "track_all_groups", "true"),
				),
			},
			{
				// Import using "user_id/true" format: track_all_groups should be set
				// and groups should be populated from the API.
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found in state: %s", resourceName)
					}
					return rs.Primary.ID + "/true", nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("expected 1 resource in imported state")
					}
					if s[0].Attributes["user_id"] == "" {
						return errors.New("user_id is empty after import")
					}
					if s[0].ID != s[0].Attributes["user_id"] {
						return errors.New("imported resource ID does not match user_id attribute")
					}
					if s[0].Attributes["track_all_groups"] != "true" {
						return errors.New("track_all_groups should be true after import with /true suffix")
					}
					if s[0].Attributes["groups.#"] == "" || s[0].Attributes["groups.#"] == "0" {
						return errors.New("groups is empty after import")
					}
					return nil
				},
			},
		},
	})
}

func TestAccResourceOktaUserGroupMemberships_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUserGroupMemberships, t.Name())
	start := mgr.GetFixtures("basic.tf", t)
	update := mgr.GetFixtures("basic_update.tf", t)
	remove := mgr.GetFixtures("basic_removal.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: start,
			},
			{
				Config: update,
			},
			{
				Config: remove,
			},
		},
	})
}
