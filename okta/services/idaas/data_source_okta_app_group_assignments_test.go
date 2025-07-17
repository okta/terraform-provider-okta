package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// TestAccDataSourceOktaAppGroupAssignments_read can flap when all the tests are
// run in harness but rarely fails running individually.
func TestAccDataSourceOktaAppGroupAssignments_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppGroupAssignments, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Check that groups exist
					resource.TestCheckResourceAttrSet("data.okta_app_group_assignments.test", "groups.#"),

					// Test group IDs match
					resource.TestCheckTypeSetElemAttrPair(
						"data.okta_app_group_assignments.test", "groups.*.id",
						"okta_group.test1", "id",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"data.okta_app_group_assignments.test", "groups.*.id",
						"okta_group.test2", "id",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"data.okta_app_group_assignments.test", "groups.*.id",
						"okta_group.test3", "id",
					),

					// Test priorities exist and match expected values
					resource.TestCheckTypeSetElemNestedAttrs("data.okta_app_group_assignments.test", "groups.*", map[string]string{
						"priority": "1",
						"profile":  "{}", // OAuth app groups have empty JSON object profiles
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.okta_app_group_assignments.test", "groups.*", map[string]string{
						"priority": "2",
						"profile":  "{}",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.okta_app_group_assignments.test", "groups.*", map[string]string{
						"priority": "3",
						"profile":  "{}",
					}),
				),
			},
		},
	})
}
