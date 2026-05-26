package idaas_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// Verify we ignore 400 "already assigned" when adding owners
func TestAccResourceOktaGroupOwners_alreadyAssignedOwner400(t *testing.T) {
	acctest.RequireSKU(t, config.SKUGovernance)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupOwners, t.Name())
	tfConfig := mgr.GetFixtures("test_already_assigned_owner.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: tfConfig,
					// If provider did not ignore the 400 already-assigned error, this step would fail
				},
			},
		})
}

// Basic create/destroy smoke test. The 404 suppression code path for
// out-of-band group deletion cannot be exercised via config manipulation
// because Terraform's dependency ordering destroys owners before the group.
func TestAccResourceOktaGroupOwners_basicCreateDestroy(t *testing.T) {
	acctest.RequireSKU(t, config.SKUGovernance)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupOwners, t.Name())
	tfConfig := mgr.GetFixtures("test_basic_create_destroy.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: tfConfig,
				},
			},
		})
}

func TestAccResourceOktaGroupOwners_groupAsOwner(t *testing.T) {
	acctest.RequireSKU(t, config.SKUGovernance)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupOwners, t.Name())
	tfConfig := mgr.GetFixtures("test_group_owner_group.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: tfConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("okta_group.child", "id"),
						resource.TestCheckTypeSetElemAttrPair("okta_group_owners.owners", "owner.*.id", "okta_group.parent", "id"),
					),
				},
			},
		})
}

func TestAccResourceOktaGroupOwners_invalidTypeMismatch(t *testing.T) {
	acctest.RequireSKU(t, config.SKUGovernance)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupOwners, t.Name())
	tfConfig := mgr.GetFixtures("test_invalid_type_mismatch.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config:      tfConfig,
					ExpectError: regexp.MustCompile(`(?i)was not found as type`),
				},
			},
		})
}

func TestAccResourceOktaGroupOwners_crudAndUpdate(t *testing.T) {
	acctest.RequireSKU(t, config.SKUGovernance)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupOwners, t.Name())
	configCreate := mgr.GetFixtures("test_resource.tf", t)
	configUpdate := mgr.GetFixtures("test_update.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
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
						resource.TestCheckResourceAttr("okta_group_owners.owners", "owner.#", "2"),
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
