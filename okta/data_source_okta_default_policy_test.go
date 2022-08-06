package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccOktaDataSourceDefaultPolicy_readPasswordPolicy(t *testing.T) {
	mgr := newFixtureManager(defaultPolicy, t.Name())
	config := testAccDataSourceDefaultPolicy(mgr.Seed, sdk.PasswordPolicyType)

	oktaResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
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

func TestAccOktaDataSourceDefaultPolicy_readIdpPolicy(t *testing.T) {
	mgr := newFixtureManager(defaultPolicy, t.Name())
	config := testAccDataSourceDefaultPolicy(mgr.Seed, sdk.IdpDiscoveryType)

	oktaResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
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
