package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaLinkDefinition(t *testing.T) {
	mgr := newFixtureManager("", linkDefinition, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", linkDefinition)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(linkDefinition, doesLinkDefinitionExist),
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
	client := sdkV2ClientForTest()
	_, response, err := client.LinkedObject.GetLinkedObjectDefinition(context.Background(), id)
	return doesResourceExist(response, err)
}

// TestAccOktaUserBaseSchemaLogin_multiple_properties This test would fail
// before the fix was implemented. The fix is to put a calling mutex on create
// and delete for the `okta_link_definition` resource. The Okta management API
// ignores parallel calls to `POST /api/v1/meta/schemas/user/linkedObjects` and
// `DELETE /api/v1/meta/schemas/user/linkedObjects` and our fix is to use a
// calling mutex in the resource to impose the equivelent of `terraform
// -parallelism=1`
func TestAccResourceOktaLinkDefinition_parallel_api_calls(t *testing.T) {
	mgr := newFixtureManager("", linkDefinition, t.Name())
	config := `
resource "okta_link_definition" "one" {
	primary_name           = "testAcc_replace_with_uuid_one"
	primary_title          = "one"
	primary_description    = "one"
	associated_name        = "testAcc_replace_with_uuid_one_a"
	associated_title       = "one_a"
	associated_description = "one_a"
}
resource "okta_link_definition" "two" {
	primary_name           = "testAcc_replace_with_uuid_two"
	primary_title          = "two"
	primary_description    = "two"
	associated_name        = "testAcc_replace_with_uuid_two_a"
	associated_title       = "two_a"
	associated_description = "two_a"
}
resource "okta_link_definition" "three" {
	primary_name           = "testAcc_replace_with_uuid_three"
	primary_title          = "three"
	primary_description    = "three"
	associated_name        = "testAcc_replace_with_uuid_three_a"
	associated_title       = "three_a"
	associated_description = "three_a"
}
resource "okta_link_definition" "four" {
	primary_name           = "testAcc_replace_with_uuid_four"
	primary_title          = "four"
	primary_description    = "four"
	associated_name        = "testAcc_replace_with_uuid_four_a"
	associated_title       = "four_a"
	associated_description = "four_a"
}
resource "okta_link_definition" "five" {
	primary_name           = "testAcc_replace_with_uuid_five"
	primary_title          = "five"
	primary_description    = "five"
	associated_name        = "testAcc_replace_with_uuid_five_a"
	associated_title       = "five_a"
	associated_description = "five_a"
}
`
	config = mgr.ConfigReplace(config)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(linkDefinition, doesLinkDefinitionExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_link_definition.one", "primary_title", "one"),
					resource.TestCheckResourceAttr("okta_link_definition.two", "primary_title", "two"),
					resource.TestCheckResourceAttr("okta_link_definition.three", "primary_title", "three"),
					resource.TestCheckResourceAttr("okta_link_definition.four", "primary_title", "four"),
					resource.TestCheckResourceAttr("okta_link_definition.five", "primary_title", "five"),
				),
			},
		},
	})
}
