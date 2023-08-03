package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccDataSourceOktaDefaultPolicy_readPasswordPolicy(t *testing.T) {
	mgr := newFixtureManager(defaultPolicy, t.Name())
	config := testAccDataSourceDefaultPolicy(mgr.Seed, sdk.PasswordPolicyType)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
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

func TestAccDataSourceOktaDefaultPolicy_readIdpPolicy(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDataSourceDefaultPolicy(ri, sdk.IdpDiscoveryType)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_default_policy.default-"+strconv.Itoa(ri), "id"),
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
