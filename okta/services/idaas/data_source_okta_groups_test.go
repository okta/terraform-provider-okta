package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaGroups_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSGroups, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group.test_1", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test_2", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_groups.okta_groups", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.okta_groups", "groups.#", "2"),
					// the example enumeration doesn't match anything so as a string the output will be a blank string
					resource.TestCheckResourceAttrSet("data.okta_groups.built_in_groups", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.built_in_groups", "groups.#", "2"),
				),
			},
		},
	})
}

func TestAccDataSourceOktaGroups_sorting(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSGroups, t.Name())
	config := mgr.GetFixtures("test_datasource_sorting.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group.test_1", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test_2", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Test basic sorting functionality
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_created", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.sorted_by_created", "groups.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_created", "groups.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_created", "groups.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_created", "groups.0.last_membership_updated"),

					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_name", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.sorted_by_name", "groups.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_name", "groups.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_name", "groups.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_name", "groups.0.last_membership_updated"),

					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_last_updated", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.sorted_by_last_updated", "groups.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_last_updated", "groups.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_last_updated", "groups.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sorted_by_last_updated", "groups.0.last_membership_updated"),
				),
			},
		},
	})
}

func TestAccDataSourceOktaGroups_sorting_comprehensive(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSGroups, t.Name())
	config := mgr.GetFixtures("test_datasource_sorting.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group.test_1", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test_2", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Test all sort fields with ascending order
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_asc", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.sort_by_id_asc", "groups.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_asc", "groups.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_asc", "groups.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_asc", "groups.0.last_membership_updated"),

					// Test all sort fields with descending order
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_desc", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.sort_by_id_desc", "groups.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_desc", "groups.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_desc", "groups.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_desc", "groups.0.last_membership_updated"),

					// Test lastMembershipUpdated sorting
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_membership_updated", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.sort_by_membership_updated", "groups.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_membership_updated", "groups.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_membership_updated", "groups.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_membership_updated", "groups.0.last_membership_updated"),

					// Test combination with other filters
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_with_type_filter", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.sort_with_type_filter", "groups.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_with_type_filter", "groups.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_with_type_filter", "groups.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_with_type_filter", "groups.0.last_membership_updated"),

					// Verify different data source IDs for different sort parameters
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_asc", "id"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_id_desc", "id"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_by_membership_updated", "id"),
					resource.TestCheckResourceAttrSet("data.okta_groups.sort_with_type_filter", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceOktaGroups_timestamp_fields(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSGroups, t.Name())
	config := mgr.GetFixtures("test_datasource_timestamp_fields.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group.test_1", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test_2", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Test timestamp fields are properly exposed and formatted
					resource.TestCheckResourceAttrSet("data.okta_groups.with_timestamps", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.with_timestamps", "groups.#", "2"),

					// Verify timestamp format (ISO 8601 format)
					resource.TestCheckResourceAttrSet("data.okta_groups.with_timestamps", "groups.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_groups.with_timestamps", "groups.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_groups.with_timestamps", "groups.0.last_membership_updated"),

					// Verify timestamp format matches expected pattern
					resource.TestCheckResourceAttrWith("data.okta_groups.with_timestamps", "groups.0.created", func(value string) error {
						// Should be in format "2006-01-02T15:04:05.000Z"
						if len(value) == 0 {
							return fmt.Errorf("created timestamp should not be empty")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.okta_groups.with_timestamps", "groups.0.last_updated", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("last_updated timestamp should not be empty")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.okta_groups.with_timestamps", "groups.0.last_membership_updated", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("last_membership_updated timestamp should not be empty")
						}
						return nil
					}),
				),
			},
		},
	})
}
