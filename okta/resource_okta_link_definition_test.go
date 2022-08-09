package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaLinkDefinition(t *testing.T) {
	mgr := newFixtureManager(linkDefinition, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", linkDefinition)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(linkDefinition, doesLinkDefinitionExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "primary_name", buildResourceName(mgr.Seed)),
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
	client := oktaClientForTest()
	_, response, err := client.LinkedObject.GetLinkedObjectDefinition(context.Background(), id)
	return doesResourceExist(response, err)
}
