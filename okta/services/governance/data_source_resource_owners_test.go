package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaResourceOwners_read(t *testing.T) {
	_ = newFixtureManager("data-sources", resources.OktaGovernanceResourceOwners, t.Name())

	_, parentOrn := discoverEntitlementBundleORN(t)

	config := fmt.Sprintf(`
data "okta_resource_owners" "test" {
  filter = "parentResourceOrn eq \"%s\""
}
`, parentOrn)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.okta_resource_owners.test", "id"),
						resource.TestCheckResourceAttrSet("data.okta_resource_owners.test", "resource_owners.#"),
					),
				},
			},
		})
}
