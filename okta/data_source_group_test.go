package okta

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGroup(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaGroup)
	config := mgr.GetFixtures("datasource.tf", ri, t)

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
