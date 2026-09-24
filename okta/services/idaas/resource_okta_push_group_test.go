package idaas_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
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

func TestAccResourceOktaPushGroup_disappears(t *testing.T) {
	if acctest.SkipVCRTest(t) {
		return
	}

	resourceName := fmt.Sprintf("%s.sample", resources.OktaIDaaSPushGroup)
	mgr := newFixtureManager("resources", resources.OktaIDaaSPushGroup, t.Name())
	config := mgr.GetFixtures("okta_push_group.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
			{
				Config:             config,
				Check:              testAccDeletePushGroupMappingOutOfBand(resourceName),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDeletePushGroupMappingOutOfBand(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		appID := rs.Primary.Attributes["app_id"]
		mappingID := rs.Primary.ID
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV6()
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		// Okta requires the mapping to be inactive before deletion.
		if _, _, err := client.GroupPushMappingAPI.UpdateGroupPushMapping(ctx, appID, mappingID).
			Body(v6okta.UpdateGroupPushMappingRequest{Status: "INACTIVE"}).Execute(); err != nil {
			return fmt.Errorf("unable to deactivate push group mapping %q: %w", mappingID, err)
		}
		if _, err := client.GroupPushMappingAPI.DeleteGroupPushMapping(ctx, appID, mappingID).
			DeleteTargetGroup(true).Execute(); err != nil {
			return fmt.Errorf("unable to delete push group mapping %q: %w", mappingID, err)
		}

		return retry.RetryContext(ctx, time.Minute, func() *retry.RetryError {
			_, response, err := client.GroupPushMappingAPI.GetGroupPushMapping(ctx, appID, mappingID).Execute()
			if response != nil && response.Response != nil && response.StatusCode == http.StatusNotFound {
				return nil
			}
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("unable to verify deletion of push group mapping %q: %w", mappingID, err))
			}

			return retry.RetryableError(fmt.Errorf("push group mapping %q still exists", mappingID))
		})
	}
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
