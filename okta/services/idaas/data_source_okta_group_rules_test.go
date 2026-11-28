package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaGroupRules_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSGroupRules, t.Name())
	config := mgr.GetFixtures("test_datasource.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group_rule.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_group_rule.test2", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_group_rules.all", "id"),
					resource.TestCheckResourceAttr("data.okta_group_rules.all", "group_rules.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_group_rules.filtered", "id"),
					resource.TestCheckResourceAttr("data.okta_group_rules.filtered", "group_rules.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_group_rules.limited", "id"),
					resource.TestCheckResourceAttr("data.okta_group_rules.limited", "group_rules.#", "2"),
					resource.TestCheckResourceAttrSet("data.okta_group_rules.with_expand", "id"),
					resource.TestCheckResourceAttr("data.okta_group_rules.with_expand", "group_rules.#", "2"),
					resource.TestCheckResourceAttr("data.okta_group_rules.with_expand", "expand", "groupIdToGroupNameMap"),
				),
			},
		},
	})
} 
