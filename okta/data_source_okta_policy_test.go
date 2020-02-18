package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourcePolicy_read(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDataSourcePolicyConfig(ri)

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

func testAccDataSourcePolicyConfig(rInt int) string {
	return fmt.Sprintf(`
data "okta_policy" "test" {
  type = "%s"
  name = "Default Policy"
}
`, passwordPolicyType)
}
