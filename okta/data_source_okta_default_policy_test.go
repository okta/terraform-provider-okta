package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccOktaDataSourceDefaultPolicy_readPasswordPolicy(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDataSourceDefaultPolicy(ri, sdk.PasswordPolicyType)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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

func TestAccOktaDataSourceDefaultPolicy_readIdpPolicy(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDataSourceDefaultPolicy(ri, sdk.IdpDiscoveryType)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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
