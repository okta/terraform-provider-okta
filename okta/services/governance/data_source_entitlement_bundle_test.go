package governance_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"testing"
)

func TestAccDataSourceOktaEntitlementBundle_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceEntitlementBundle, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_entitlement_bundle.test", "id"),
					resource.TestCheckResourceAttr("data.okta_entitlement_bundle.test", "name", "entitlement bundle data source test"),
					resource.TestCheckResourceAttr("data.okta_entitlement_bundle.test", "target.type", "APPLICATION"),
				),
			},
		},
	})
}
