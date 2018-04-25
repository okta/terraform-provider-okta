package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceDefaultPolicy(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("data.okta_default_policies.default-"+strconv.Itoa(ri), "id"),
				),
			},
		},
	})
}

func testAccDataSourceDefaultPolicyConfig(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "PASSWORD"
}
`, rInt)
}
