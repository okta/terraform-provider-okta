package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaGroup_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", group, t.Name())
	groupCreate := mgr.GetFixtures("okta_group.tf", t)
	config := mgr.GetFixtures("datasource.tf", t)
	configInvalid := mgr.GetFixtures("datasource_not_found.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: groupCreate,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_group.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_group.test", "type"),
					resource.TestCheckResourceAttr("data.okta_group.test", "users.#", "1"),
				),
			},
			{
				Config:      configInvalid,
				ExpectError: regexp.MustCompile(`\bdoes not exist`),
			},
		},
	})
}

// TestAccDataSourceOktaGroup_read_multiple_groups Checks of the group data
// source surfaces the correct error when more than one group matches the name
// argument.
func TestAccDataSourceOktaGroup_read_multiple_groups(t *testing.T) {
	resourceName := fmt.Sprintf("data.%s.test", group)
	mgr := newFixtureManager("data-sources", group, t.Name())

	baseConfig := `
resource "okta_group" "test_1" {
  name        = "testAcc_replace_with_uuid_MORE"
  description = "testing, more"
}
resource "okta_group" "test_2" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}`

	step1Config := `
# error, name is substring and matches both groups
data "okta_group" "test" {
  name          = "testAcc"
  depends_on = [okta_group.test_1, okta_group.test_2]
}`

	step2Config := `
# ok, name is substring and matches both groups, but is exact match for one group
data "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  depends_on = [okta_group.test_1, okta_group.test_2]
}`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", step1Config, baseConfig)),
				// NOTE: there might be dangling test groups on the org starting
				// with "testAcc" so just make sure the error message is correct
				// besides the count of groups
				ExpectError: regexp.MustCompile(`group starting with name "testAcc" matches (\d+) groups, select a more precise name parameter`),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", step2Config, baseConfig)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "testing, testing"),
				),
			},
		},
	})
}
