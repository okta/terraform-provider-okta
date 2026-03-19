package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaGroupRole_admin_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRole)
	resourceName2 := fmt.Sprintf("%s.test_app", resources.OktaIDaaSGroupRole)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRole, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	groupTarget := mgr.GetFixtures("group_targets.tf", t)
	groupTargetsUpdated := mgr.GetFixtures("group_targets_updated.tf", t)
	groupTargetsRemoved := mgr.GetFixtures("group_targets_removed.tf", t)

	// NOTE this test doesn't always pass
	// "The role specified is already assigned to the user."
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "READ_ONLY_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "target_app_list.#", "0"),
				),
			},
			{
				Config: groupTarget,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "HELP_DESK_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "target_app_list.#", "1"),
				),
			},
			{
				Config: groupTargetsUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "HELP_DESK_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "target_app_list.#", "1"),
				),
			},
			{
				Config: groupTargetsRemoved,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "HELP_DESK_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "target_app_list.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceOktaGroupRole_custom_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRole)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRole, t.Name())
	config := mgr.GetFixtures("custom.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "CUSTOM"),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_set_id"),
				),
			},
		},
	})
}
