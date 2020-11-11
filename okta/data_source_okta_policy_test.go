package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func TestAccOktaDataSourcePolicy_read(t *testing.T) {
	config := testAccDataSourcePolicyConfig()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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
