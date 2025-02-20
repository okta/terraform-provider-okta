package idaas_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccDataSourceOktaDefaultPolicy_readPasswordPolicy(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSDefaultPolicy, t.Name())
	config := testAccDataSourceDefaultPolicy(mgr.Seed, sdk.PasswordPolicyType)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_default_policy.default-"+strconv.Itoa(mgr.Seed), "id"),
				),
			},
		},
	})
}

func TestAccDataSourceOktaDefaultPolicy_readIdpPolicy(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSDefaultPolicy, t.Name())
	config := testAccDataSourceDefaultPolicy(mgr.Seed, sdk.IdpDiscoveryType)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_default_policy.default-"+strconv.Itoa(mgr.Seed), "id"),
				),
			},
		},
	})
}

func testAccDataSourceDefaultPolicy(rInt int, policy string) string {
	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
  type = "%s"
}
`, rInt, policy)
}
