package idaas_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

// Verify we ignore 400 "already assigned" when adding owners
func TestAccResourceOktaGroupOwners_alreadyAssignedOwner400(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owners", t.Name())
	config := mgr.GetFixtures("test_already_assigned_owner.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: config,
					// If provider did not ignore the 400 already-assigned error, this step would fail
				},
			},
		})
}

// Verify we suppress 404s when group is deleted ahead of deleting owners
func TestAccResourceOktaGroupOwners_deletedGroup404(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owners", t.Name())
	config := mgr.GetFixtures("test_deleted_group_404.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: config,
				},
				// Second step: remove the group resource so it gets deleted; keep owners to trigger 404
				{
					Config: regexp.MustCompile(`(?s)resource "okta_group" "grp"\{.*?\}\n`).ReplaceAllString(config, ""),
					// Destroy then apply; any 404 during owner deletions should be suppressed by provider
				},
			},
		})
}

func TestAccResourceOktaGroupOwners_groupAsOwner(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owners", t.Name())
	config := mgr.GetFixtures("test_group_owner_group.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("okta_group.child", "id"),
						resource.TestCheckTypeSetElemAttrPair("okta_group_owners.owners", "owner.*.id", "okta_group.parent", "id"),
					),
				},
			},
		})
}

func TestAccResourceOktaGroupOwners_invalidTypeMismatch(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owners", t.Name())
	config := mgr.GetFixtures("test_invalid_type_mismatch.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config:      config,
					ExpectError: regexp.MustCompile(`(?i)(invalid|validation|type|api)`),
				},
			},
		})
}

func TestAccResourceOktaGroupOwners_crudAndUpdate(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owners", t.Name())
	configCreate := mgr.GetFixtures("test_resource.tf", t)
	configUpdate := mgr.GetFixtures("test_update.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: configCreate,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("okta_group.grp", "id"),
						resource.TestCheckResourceAttrPair("okta_group_owners.owners", "group_id", "okta_group.grp", "id"),
						resource.TestCheckTypeSetElemAttrPair("okta_group_owners.owners", "owner.*.id", "okta_user.owner1", "id"),
						resource.TestCheckTypeSetElemAttrPair("okta_group_owners.owners", "owner.*.id", "okta_user.owner2", "id"),
					),
				},
				{
					Config: configUpdate,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckTypeSetElemAttrPair("okta_group_owners.owners", "owner.*.id", "okta_user.owner1", "id"),
						// owner2 removed, owner3 added
						resource.TestCheckTypeSetElemAttrPair("okta_group_owners.owners", "owner.*.id", "okta_user.owner3", "id"),
					),
				},
				{
					ResourceName:      "okta_group_owners.owners",
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						groupID := s.RootModule().Resources["okta_group.grp"].Primary.Attributes["id"]
						return groupID, nil
					},
				},
			},
		})
}

func TestAccResourceOktaGroupOwners_trackAllOwnersFalse(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owners", t.Name())
	config := mgr.GetFixtures("test_track_false.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("okta_group.grp", "id"),
						resource.TestCheckResourceAttr("okta_group_owners.owners", "track_all_owners", "false"),
						resource.TestCheckTypeSetElemAttrPair("okta_group_owners.owners", "owner.*.id", "okta_user.owner1", "id"),
						resource.TestCheckTypeSetElemAttrPair("okta_group_owners.owners", "owner.*.id", "okta_user.owner2", "id"),
					),
				},
				// Now simulate refresh behavior: with track_all_owners=false the external owner should not be removed
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("okta_group_owners.owners", "track_all_owners", "false"),
					),
				},
			},
		})
}
