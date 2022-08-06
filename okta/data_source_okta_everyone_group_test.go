package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceEveryoneGroup_read(t *testing.T) {
	mgr := newFixtureManager(groupEveryone, t.Name())
	config := testAccDataSourceEveryoneGroupConfig(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_everyone_group.everyone-"+strconv.Itoa(mgr.Seed), "id"),
				),
			},
		},
	})
}

func testAccDataSourceEveryoneGroupConfig(rInt int) string {
	return fmt.Sprintf(`
data "okta_everyone_group" "everyone-%d" {}
`, rInt)
}
