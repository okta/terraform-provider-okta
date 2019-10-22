package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourceDefaultPolicy_read(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDataSourceDefaultPolicyConfig(ri)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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

func testAccDataSourceDefaultPolicyConfig(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
  type = "%s"
}
`, rInt, passwordPolicyType)
}
