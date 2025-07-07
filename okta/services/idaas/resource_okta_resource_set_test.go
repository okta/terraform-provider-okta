package idaas_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

// TODO: V6 BREAKING CHANGE - These tests will need to be updated when okta_resource_set is split
// into okta_resource_set and okta_resource_set_resource in V6.
// We will need to ensure that the tests are updated to test both okta_resource_set and okta_resource_set_resources
func TestAccResourceOktaResourceSet_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSResourceSet, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSResourceSet)
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSResourceSet, doesResourceSetExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing"),
						resource.TestCheckResourceAttr(resourceName, "resources.#", "3"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing updated"),
						resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
					),
				},
			},
		})
}

// TestAccResourceOktaResourceSet_Issue1097_Pagination deals with resolving a
// pagination bug with more than 100 resources in the set
// https://github.com/okta/terraform-provider-okta/issues/1097
//
// OKTA_ALLOW_LONG_RUNNING_ACC_TEST=true TF_ACC=1 \
// go test -tags unit -mod=readonly -test.v -run ^TestAccResourceOktaResourceSet_Issue1097_Pagination$ ./okta 2>&1
func TestAccResourceOktaResourceSet_Issue1097_Pagination(t *testing.T) {
	if !allowLongRunningACCTest(t) {
		t.SkipNow()
	}

	orgName := os.Getenv("OKTA_ORG_NAME")
	baseUrl := os.Getenv("OKTA_BASE_URL")
	config := fmt.Sprintf(`
		resource "okta_group" "testing" {
			count = 201
			name = "group_replace_with_uuid_${count.index}"
		}

		resource "okta_resource_set" "test" {
			label       = "testAcc_replace_with_uuid"
			description = "set of resources"

			resources = [
				for group in okta_group.testing :
					"https://%s.%s/api/v1/groups/${group.id}"
			]
		}`, orgName, baseUrl)
	mgr := newFixtureManager("resources", resources.OktaIDaaSResourceSet, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSResourceSet)
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSResourceSet, doesResourceSetExist),
			Steps: []resource.TestStep{
				{
					Config: mgr.ConfigReplace(config),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
						// NOTE: before bug fix test would error out on having a
						// detected change of extra items in the resources list
						// beyond the first 100 resources.
						//
						// Plan: 0 to add, 1 to change, 0 to destroy.
						resource.TestCheckResourceAttr(resourceName, "resources.#", "201"),
					),
				},
			},
		})
}

// TestAccResourceOktaResourceSet_Issue_1735_drift_detection
// This test demonstrates that resource okta_resource_set implements proper drift detection.
func TestAccResourceOktaResourceSet_Issue_1735_drift_detection(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSResourceSet, t.Name())
	resourceSet := fmt.Sprintf("%s.test", resources.OktaIDaaSResourceSet)

	baseConfig := `
variable "hostname" {
  type = string
}`

	step1Config := `
resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  resources = [
    "https://${var.hostname}/api/v1/users",
    "https://${var.hostname}/api/v1/groups"
  ]
}`

	step2Config := `
resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  resources = [
    "https://${var.hostname}/api/v1/users",
    "https://${var.hostname}/api/v1/groups",
    "https://${var.hostname}/api/v1/apps",
  ]
}`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step1Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceSet, "resources.#", "2"),
					resource.TestCheckResourceAttr(resourceSet, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(resourceSet, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step1Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceSet, "resources.#", "2"),
					resource.TestCheckResourceAttr(resourceSet, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(resourceSet, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
					// This mimics adding the apps resource to the resource set
					// outside of Terraform.  In this case doing so with a
					// direct API call via the test harness which is equivalent
					// to "Click Ops"
					clickOpsAddResourceToResourceSet(resources.OktaIDaaSResourceSet+".test", fmt.Sprintf("https://%s/api/v1/apps", os.Getenv("TF_VAR_hostname"))),

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
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step2Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceSet, "resources.#", "3"),
					resource.TestCheckResourceAttr(resourceSet, "resources.0", fmt.Sprintf("https://%s/api/v1/apps", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(resourceSet, "resources.1", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(resourceSet, "resources.2", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step1Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceSet, "resources.#", "2"),
					resource.TestCheckResourceAttr(resourceSet, "resources.0", fmt.Sprintf("https://%s/api/v1/groups", os.Getenv("TF_VAR_hostname"))),
					resource.TestCheckResourceAttr(resourceSet, "resources.1", fmt.Sprintf("https://%s/api/v1/users", os.Getenv("TF_VAR_hostname"))),
				),
			},
		},
	})
}

func TestAccResourceOktaResourceSet_Issue_1991_orn_handling(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSResourceSet, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSResourceSet)

	baseConfig := `
variable "id" {
  type = string
}
`
	step1Config := `
resource "okta_resource_set" "test" {
  label       = "testAcc_replace_1991"
  description = "testing, testing"
  resources_orn = [
    "orn:okta:directory:${var.id}:users",
    "orn:okta:directory:${var.id}:groups"
  ]
}`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step1Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "resources_orn.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "resources_orn.0", fmt.Sprintf("orn:okta:directory:%s:groups", os.Getenv("TF_VAR_orgID"))),
					resource.TestCheckResourceAttr(resourceName, "resources_orn.1", fmt.Sprintf("orn:okta:directory:%s:users", os.Getenv("TF_VAR_orgID"))),
				),
			},
		},
	})
}

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

func doesResourceSetExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
	_, response, err := client.GetResourceSet(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}
