package idaas_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

// TestAccResourceOktaResourceSet_crud tests basic CRUD operations for the resource set.
// This ensures create, read, update, and delete operations work as expected.
func TestAccResourceOktaResourceSet_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSResourceSet, t.Name())
	stateAddress := fmt.Sprintf("%s.test", resources.OktaIDaaSResourceSet)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSResourceSet, doesResourceSetExist),
			Steps: []resource.TestStep{
				{
					Config: mgr.GetFixtures("test_basic.tf", t),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(stateAddress, "label", acctest.BuildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(stateAddress, "description", "testing, testing"),
						resource.TestCheckResourceAttr(stateAddress, "resources.#", "3"),
					),
				},
				{
					Config: mgr.GetFixtures("test_basic_updated.tf", t),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(stateAddress, "label", fmt.Sprintf("%s updated", acctest.BuildResourceName(mgr.Seed))),
						resource.TestCheckResourceAttr(stateAddress, "description", "testing, testing updated"),
						resource.TestCheckResourceAttr(stateAddress, "resources.#", "2"),
					),
				},
			},
		})
}

// TestAccResourceOktaResourceSet_Issue1097_Pagination tests the fix for
// https://github.com/okta/terraform-provider-okta/issues/1097
// where pagination would fail with more than 200 resources.
//
// Uses 201 resources to specifically test handling across Okta's 200-item
// page boundary. The issue manifested as:
// - Resources beyond the first 200 would be lost
// - State refresh would fail to capture all resources
// - Plan would show phantom changes
func TestAccResourceOktaResourceSet_Issue1097_Pagination(t *testing.T) {
	if !allowLongRunningACCTest(t) {
		t.SkipNow()
	}

	mgr := newFixtureManager("resources", resources.OktaIDaaSResourceSet, t.Name())
	stateAddress := fmt.Sprintf("%s.test", resources.OktaIDaaSResourceSet)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSResourceSet, doesResourceSetExist),
			Steps: []resource.TestStep{
				{
					Config: mgr.GetFixtures("test_pagination.tf", t),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(stateAddress, "resources.#", "201"),
					),
				},
			},
		})
}

// TestAccResourceOktaResourceSet_Issue_1735_drift_detection verifies that
// the resource properly detects and handles changes made outside of Terraform.
// https://github.com/okta/terraform-provider-okta/issues/1735
//
// This test simulates a scenario where:
// 1. Resources are created via Terraform
// 2. Additional resources are added outside of Terraform ("click ops")
// 3. Terraform detects these changes and handles them appropriately
// 4. Changes can be reconciled back to the desired state
func TestAccResourceOktaResourceSet_Issue_1735_drift_detection(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSResourceSet, t.Name())
	stateAddress := fmt.Sprintf("%s.test", resources.OktaIDaaSResourceSet)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("test_drift_detection.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources.#", "2"),
					resource.TestCheckResourceAttr(stateAddress, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
			{
				Config: mgr.GetFixtures("test_drift_detection.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources.#", "2"),
					resource.TestCheckResourceAttr(stateAddress, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
					// This mimics adding the apps resource to the resource set
					// outside of Terraform (e.g., via UI or API directly).
					// This simulates "click ops" or manual changes that Terraform
					// needs to detect and handle.
					clickOpsAddResourceToResourceSet(stateAddress, fmt.Sprintf("https://%s/api/v1/apps", os.Getenv("TF_VAR_hostname"))),
				),
				ExpectNonEmptyPlan: true, // Plan will show difference due to external modification
			},
			{
				Config: mgr.GetFixtures("test_drift_detection_updated.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources.#", "3"),
					resource.TestCheckResourceAttr(stateAddress, "resources.0", fmt.Sprintf("https://%s/api/v1/apps", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.1", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.2", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
			{
				Config: mgr.GetFixtures("test_drift_detection.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources.#", "2"),
					resource.TestCheckResourceAttr(stateAddress, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
		},
	})
}

// TestAccResourceOktaResourceSet_Issue_1991_orn_handling verifies support for
// Okta Resource Names (ORNs) in resource sets. This ensures compatibility with
// Okta's newer ORN-based resource addressing scheme.
// https://github.com/okta/terraform-provider-okta/issues/1991
func TestAccResourceOktaResourceSet_Issue_1991_orn_handling(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSResourceSet, t.Name())
	stateAddress := fmt.Sprintf("%s.test", resources.OktaIDaaSResourceSet)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("test_orn_handling.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources_orn.#", "2"),
					// Use regex patterns to validate ORN format without being tied to specific org ID
					func(s *terraform.State) error {
						// Get the actual ORN values from the state
						actualORN0 := s.RootModule().Resources[stateAddress].Primary.Attributes["resources_orn.0"]
						actualORN1 := s.RootModule().Resources[stateAddress].Primary.Attributes["resources_orn.1"]

						// Regex pattern for ORN format: orn:okta:directory:<org-id>:<resource-type>
						// This accepts any org ID (including empty) and validates the overall structure
						ornPattern := `^orn:okta:directory:[^:]*:(groups|users)$`

						// Check the first ORN (should be groups)
						if !regexp.MustCompile(ornPattern).MatchString(actualORN0) {
							return fmt.Errorf("resources_orn.0: expected format 'orn:okta:directory:<org-id>:groups', got %q", actualORN0)
						}

						// Check the second ORN (should be users)
						if !regexp.MustCompile(ornPattern).MatchString(actualORN1) {
							return fmt.Errorf("resources_orn.1: expected format 'orn:okta:directory:<org-id>:users', got %q", actualORN1)
						}

						// Additional validation: ensure we have exactly one groups and one users ORN
						groupsCount := 0
						usersCount := 0
						for i := 0; i < 2; i++ {
							orn := s.RootModule().Resources[stateAddress].Primary.Attributes[fmt.Sprintf("resources_orn.%d", i)]
							if strings.HasSuffix(orn, ":groups") {
								groupsCount++
							} else if strings.HasSuffix(orn, ":users") {
								usersCount++
							}
						}

						if groupsCount != 1 || usersCount != 1 {
							return fmt.Errorf("expected exactly 1 groups ORN and 1 users ORN, got %d groups and %d users", groupsCount, usersCount)
						}

						return nil
					},
				),
			},
		},
	})
}

// clickOpsAddResourceToResourceSet simulates adding a resource to a resource set
// outside of Terraform (e.g., through the Okta UI or direct API calls).
// This helper is used to test drift detection and state reconciliation.
func clickOpsAddResourceToResourceSet(resourceSet, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource set not found: %s", resources.OktaIDaaSResourceSet)
		resourceSetRS, ok := s.RootModule().Resources[resourceSet]
		if !ok {
			return missingErr
		}
		resourceSetID := resourceSetRS.Primary.Attributes["id"]

		client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
		patch := sdk.PatchResourceSet{Additions: []string{resourceName}}
		_, _, err := client.PatchResourceSet(context.Background(), resourceSetID, patch)
		if err != nil {
			return fmt.Errorf("API: unable to patch resource %q with addition %q, err: %+v", resourceSetID, resourceName, err)
		}

		return nil
	}
}

// doesResourceSetExist verifies whether a resource set exists in Okta.
// Used by the test framework to validate resource creation/deletion.
func doesResourceSetExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
	_, response, err := client.GetResourceSet(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}
