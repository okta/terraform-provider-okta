package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaGroupAdminRole_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", groupRole)
	mgr := newFixtureManager(groupRole)
	config := mgr.GetFixtures("basic.tf", ri, t)
	groupTarget := mgr.GetFixtures("group_targets.tf", ri, t)
	groupTargetsUpdated := mgr.GetFixtures("group_targets_updated.tf", ri, t)
	groupTargetsRemoved := mgr.GetFixtures("group_targets_removed.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(oktaGroup, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(resourceName, "role_type", "READ_ONLY_ADMIN"),
			},
			{
				Config: groupTarget,
				Check:  resource.TestCheckResourceAttrSet(resourceName, "group_target_list"),
			},
			{
				Config: groupTargetsUpdated,
				Check:  resource.TestCheckResourceAttrSet(resourceName, "group_target_list"),
			},
			{
				Config: groupTargetsRemoved,
				Check:  resource.TestCheckResourceAttr(resourceName, "role_type", "HELP_DESK_ADMIN"),
			},
		},
	})
}
