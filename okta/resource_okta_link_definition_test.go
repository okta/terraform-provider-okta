package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func sweepLinkDefinitions(client *testClient) error {
	var errorList []error
	linkedObjects, _, err := getOktaClientFromMetadata(testAccProvider.Meta()).LinkedObject.ListLinkedObjectDefinitions(context.Background())
	if err != nil {
		return err
	}
	for _, object := range linkedObjects {
		if strings.HasPrefix(object.Primary.Name, testResourcePrefix) {
			if _, err := getOktaClientFromMetadata(testAccProvider.Meta()).LinkedObject.DeleteLinkedObjectDefinition(context.Background(), object.Primary.Name); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func TestAccOktaLinkDefinition(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(linkDefinition)
	config := mgr.GetFixtures("basic.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", linkDefinition)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(linkDefinition, doesLinkDefinitionExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "primary_name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "primary_title", "Manager"),
					resource.TestCheckResourceAttr(resourceName, "primary_description", "Manager link property"),
					resource.TestCheckResourceAttr(resourceName, "associated_name", "testAcc_subordinate"),
					resource.TestCheckResourceAttr(resourceName, "associated_title", "Subordinate"),
					resource.TestCheckResourceAttr(resourceName, "associated_description", "Subordinate link property"),
				),
			},
		},
	})
}

func doesLinkDefinitionExist(id string) (bool, error) {
	_, response, err := getOktaClientFromMetadata(testAccProvider.Meta()).LinkedObject.GetLinkedObjectDefinition(context.Background(), id)
	return doesResourceExist(response, err)
}
