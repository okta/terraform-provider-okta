package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaPushGroup_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.sample", resources.OktaIDaaSPushGroup)
	mgr := newFixtureManager("resources", resources.OktaIDaaSPushGroup, t.Name())
	config := mgr.GetFixtures("okta_push_group.tf", t)
	updatedConfig := mgr.GetFixtures("okta_push_group_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_target_group_on_destroy", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "source_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_group_id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_target_group_on_destroy", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "source_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_group_id"),
				),
			},
		},
	})
}

func TestAccResourceOktaPushGroup_ad(t *testing.T) {
	resourceName := fmt.Sprintf("%s.sample", resources.OktaIDaaSPushGroup)
	mgr := newFixtureManager("resources", resources.OktaIDaaSPushGroup, t.Name())
	config := mgr.GetFixtures("okta_push_group_ad.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_target_group_on_destroy", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "source_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_group_id"),
					resource.TestCheckResourceAttr(resourceName, "app_config.type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceName, "app_config.distinguished_name", "CN=Test,OU=Groups,DC=example,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "app_config.group_scope", "DOMAIN_LOCAL"),
					resource.TestCheckResourceAttr(resourceName, "app_config.group_type", "SECURITY"),
					resource.TestCheckResourceAttr(resourceName, "app_config.sam_account_name", "something"),
				),
			},
		},
	})
}
