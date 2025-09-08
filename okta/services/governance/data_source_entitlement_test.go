package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEntitlement_read(t *testing.T) {
	t.Skip("Skipping Entitlement tests")
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceEntitlement, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_entitlement.test", "id"),
					resource.TestCheckResourceAttr("data.okta_entitlement.test", "name", "Entitlement Bundle"),
					resource.TestCheckResourceAttr("data.okta_entitlement.test", "external_value", "Entitlement Bundle"),
					resource.TestCheckResourceAttr("data.okta_entitlement.test", "parent.type", "APPLICATION"),
				),
			},
		},
	})
}
