package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaBehaviorNewAsn_basic(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_behavior_new_asn", t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := "data.okta_behavior_new_asn.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", "okta_behavior_new_asn.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc_"+mgr.SeedStr()),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated"),
					resource.TestCheckResourceAttr(resourceName, "settings.max_events_used_for_evaluation", "20"),
					resource.TestCheckResourceAttr(resourceName, "settings.min_events_needed_for_evaluation", "0"),
				),
			},
		},
	})
}
