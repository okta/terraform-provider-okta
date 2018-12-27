package okta

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOktaGroupsCreate(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := buildResourceFQN(oktaGroup, ri)
	mgr := newFixtureManager("okta_group")
	config := mgr.GetFixtures("okta_group.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_group_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testgroupdifferent")),
			},
		},
	})
}
