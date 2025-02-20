package idaas_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEveryoneGroup_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSGroupEveryone, t.Name())
	config := testAccDataSourceEveryoneGroupConfig(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
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
