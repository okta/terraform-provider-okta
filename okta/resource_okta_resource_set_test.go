package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func sweepResourceSets(client *testClient) error {
	var errorList []error
	resourceSets, _, err := client.apiSupplement.ListResourceSets(context.Background())
	if err != nil {
		return err
	}
	for _, b := range resourceSets.ResourceSets {
		if !strings.HasPrefix(b.Label, "testAcc_") {
			if _, err := client.apiSupplement.DeleteResourceSet(context.Background(), b.Id); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func TestAccOktaResourceSet(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(resourceSet)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", resourceSet)
	resource.Test(
		t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
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
