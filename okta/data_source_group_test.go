package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGroup(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDataSourceGroupConfig(ri)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_group.testAcc_"+strconv.Itoa(ri), "id"),
					resource.TestCheckResourceAttrSet("okta_group.testAcc_"+strconv.Itoa(ri), "id"),
				),
			},
		},
	})
}

func testAccDataSourceGroupConfig(ri int) string {
	return fmt.Sprintf(`
resource "okta_group" "testAcc_%[1]d" {
	name        = "something new"
	description = "testing, testing"
}

data "okta_group" "testAcc_%[1]d" {
	name = "${okta_group.testAcc_%[1]d.name}"
}
`, ri)
}
