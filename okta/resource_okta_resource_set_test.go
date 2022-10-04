package okta

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaResourceSet(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(resourceSet)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", resourceSet)
	os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL")))
	resource.Test(
		t, resource.TestCase{
			PreCheck:          testAccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createCheckResourceDestroy(resourceSet, doesResourceSetExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing"),
						resource.TestCheckResourceAttr(resourceName, "resources.#", "3"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
						resource.TestCheckResourceAttr(resourceName, "description", "testing, testing updated"),
						resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
					),
				},
			},
		})
}

func doesResourceSetExist(id string) (bool, error) {
	_, response, err := getSupplementFromMetadata(testAccProvider.Meta()).GetResourceSet(context.Background(), id)
	return doesResourceExist(response, err)
}

// TestAccResrouceOktaResourceSet_Issue1097_Pagination deals with resolving a
// pagination bug with more than 100 resources in the set
// https://github.com/okta/terraform-provider-okta/issues/1097
//
// OKTA_ALLOW_LONG_RUNNING_ACC_TEST=true TF_ACC=1 \
// go test -tags unit -mod=readonly -test.v -run ^TestAccResrouceOktaResourceSet_Issue1097_Pagination$ ./okta 2>&1
func TestAccResrouceOktaResourceSet_Issue1097_Pagination(t *testing.T) {
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
	ri := acctest.RandInt()
	mgr := newFixtureManager(resourceSet)
	resourceName := fmt.Sprintf("%s.test", resourceSet)
	resource.Test(
		t, resource.TestCase{
			PreCheck:          testAccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createCheckResourceDestroy(resourceSet, doesResourceSetExist),
			Steps: []resource.TestStep{
				{
					Config: mgr.ConfigReplace(config, ri),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
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
