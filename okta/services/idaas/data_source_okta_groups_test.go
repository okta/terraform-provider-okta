package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaGroups_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSGroups, t.Name())
	//groups := mgr.GetFixtures("okta_groups.tf", t)
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
