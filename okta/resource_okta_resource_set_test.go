package okta

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaResourceSet_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resourceSet, t.Name())
	stateAddress := fmt.Sprintf("%s.test", resourceSet)

	TestAccResourceSetBasic := `
	variable "hostname" { type = string }
	resource "okta_resource_set" "test" {
		label       = "testAcc_replace_with_uuid"
		description = "testing, testing"
		resources = [
			"https://${var.hostname}/api/v1/users",
			"https://${var.hostname}/api/v1/apps",
			"https://${var.hostname}/api/v1/groups"
		]
	}
	`
	TestAccResourceSetUpdated := `
	variable "hostname" { type = string }
	resource "okta_resource_set" "test" {
		label       = "testAcc_replace_with_uuid updated"
		description = "testing, testing updated"
		resources = [
			"https://${var.hostname}/api/v1/users",
			"https://${var.hostname}/api/v1/apps",
		]
	}
	`

	oktaResourceTest(
		t, resource.TestCase{
			PreCheck:          testAccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      checkResourceDestroy(resourceSet, doesResourceSetExist),
			Steps: []resource.TestStep{
				{
					Config: mgr.ConfigReplace(TestAccResourceSetBasic),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(stateAddress, "label", buildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(stateAddress, "description", "testing, testing"),
						resource.TestCheckResourceAttr(stateAddress, "resources.#", "3"),
					),
				},
				{
					Config: mgr.ConfigReplace(TestAccResourceSetUpdated),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(stateAddress, "label", fmt.Sprintf("%s updated", buildResourceName(mgr.Seed))),
						resource.TestCheckResourceAttr(stateAddress, "description", "testing, testing updated"),
						resource.TestCheckResourceAttr(stateAddress, "resources.#", "2"),
					),
				},
			},
		})
}

func TestAccResourceOktaResourceSet_Issue1097_Pagination(t *testing.T) {
	// TestAccResourceOktaResourceSet_Issue1097_Pagination deals with resolving a
	// pagination bug with more than 100 resources in the set
	// https://github.com/okta/terraform-provider-okta/issues/1097
	//
	// OKTA_ALLOW_LONG_RUNNING_ACC_TEST=true TF_ACC=1 \
	// go test -tags unit -mod=readonly -test.v -run ^TestAccResourceOktaResourceSet_Issue1097_Pagination$ ./okta 2>&1
	if !allowLongRunningACCTest(t) {
		t.SkipNow()
	}

	mgr := newFixtureManager("resources", resourceSet, t.Name())
	stateAddress := fmt.Sprintf("%s.test", resourceSet)

	step1Config := `
	variable "hostname" { type = string }
	resource "okta_group" "test" {
		count = 101
		name = "testAcc_${count.index}"
	}
	resource "okta_resource_set" "test" {
		label       = "testAcc_replace_with_uuid"
		description = "testAcc_replace_with_uuid"
		resources = [ for g in okta_group.test[*] : "https://${var.hostname}/api/v1/groups/${g.id}" ]
	}
	`

	oktaResourceTest(
		t, resource.TestCase{
			PreCheck:          testAccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      checkResourceDestroy(resourceSet, doesResourceSetExist),
			Steps: []resource.TestStep{
				{
					Config: mgr.ConfigReplace(step1Config),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(stateAddress, "label", buildResourceName(mgr.Seed)),
						// NOTE: before bug fix test would error out on having a
						// detected change of extra items in the resources list
						// beyond the first 100 resources.
						//
						// Plan: 0 to add, 1 to change, 0 to destroy.
						resource.TestCheckResourceAttr(stateAddress, "resources.#", "101"),
					),
				},
			},
		})
}

func TestAccResourceOktaResourceSet_Issue_1735_drift_detection(t *testing.T) {
	// TestAccResourceOktaResourceSet_Issue_1735_drift_detection
	// This test demonstrates that resource okta_resource_set implements proper drift detection.
	mgr := newFixtureManager("resources", resourceSet, t.Name())
	stateAddress := fmt.Sprintf("%s.test", resourceSet)

	step1Config := `
	variable "hostname" { type = string }
	resource "okta_resource_set" "test" {
		label       = "testAcc_replace_with_uuid"
		description = "testing, testing"
		resources = [
			"https://${var.hostname}/api/v1/users",
			"https://${var.hostname}/api/v1/groups",
		]
	}
  	`

	step2Config := `
	variable "hostname" { type = string }
	resource "okta_resource_set" "test" {
		label       = "testAcc_replace_with_uuid"
		description = "testing, testing"
		resources = [
			"https://${var.hostname}/api/v1/users",
			"https://${var.hostname}/api/v1/groups",
			"https://${var.hostname}/api/v1/apps",
		]
	}
	`

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		CheckDestroy:      nil,
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(step1Config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources.#", "2"),
					resource.TestCheckResourceAttr(stateAddress, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
			{
				Config: mgr.ConfigReplace(step1Config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources.#", "2"),
					resource.TestCheckResourceAttr(stateAddress, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
					// This mimics adding the apps resource to the resource set
					// outside of Terraform.  In this case doing so with a
					// direct API call via the test harness which is equivalent
					// to "Click Ops"
					clickOpsAddResourceToResourceSet(stateAddress, fmt.Sprintf("https://%s/api/v1/apps", os.Getenv("TF_VAR_hostname"))),

					// NOTE: after these checks run the terraform test runner
					// will do a refresh and catch that apps resource has been
					// added to the resource set outside of the terraform config
					// and emit a non-empty plan
				),

				// side effect of the TF test runner is expecting a non-empty
				// plan is treated as an apply accept and updates the resources
				// on the resource set to the local state
				ExpectNonEmptyPlan: true,
			},
			{
				Config: mgr.ConfigReplace(step2Config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources.#", "3"),
					resource.TestCheckResourceAttr(stateAddress, "resources.0", fmt.Sprintf("https://%s/api/v1/apps", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.1", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.2", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
			{
				Config: mgr.ConfigReplace(step1Config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateAddress, "resources.#", "2"),
					resource.TestCheckResourceAttr(stateAddress, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(stateAddress, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
		},
	})
}

func clickOpsAddResourceToResourceSet(resourceSet, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource set not found: %s", resourceSet)
		resourceSetRS, ok := s.RootModule().Resources[resourceSet]
		if !ok {
			return missingErr
		}
		resourceSetID := resourceSetRS.Primary.Attributes["id"]

		client := sdkSupplementClientForTest()
		patch := sdk.PatchResourceSet{Additions: []string{resourceName}}
		_, _, err := client.PatchResourceSet(context.Background(), resourceSetID, patch)
		if err != nil {
			return fmt.Errorf("API: unable to patch resource %q with addition %q, err: %+v", resourceSetID, resourceName, err)
		}

		return nil
	}
}

func doesResourceSetExist(id string) (bool, error) {
	client := sdkSupplementClientForTest()
	_, response, err := client.GetResourceSet(context.Background(), id)
	return doesResourceExist(response, err)
}
