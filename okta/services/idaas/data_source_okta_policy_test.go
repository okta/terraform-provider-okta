package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccDataSourceOktaPolicy_read(t *testing.T) {
	config := testAccDataSourcePolicyConfig()

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_policy.test", "id"),
				),
			},
		},
	})
}

func testAccDataSourcePolicyConfig() string {
	return fmt.Sprintf(`
data "okta_policy" "test" {
  type = "%s"
  name = "Default Policy"
}
`, sdk.PasswordPolicyType)
}
