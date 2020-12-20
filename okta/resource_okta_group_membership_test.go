package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOktaGroupMembership_crud(t *testing.T) {
	ri := acctest.RandInt()

	mgr := newFixtureManager(oktaGroupMembership)
	config := mgr.GetFixtures("okta_group_membership.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_group_membership_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: updatedConfig,
			},
		},
	})
}
