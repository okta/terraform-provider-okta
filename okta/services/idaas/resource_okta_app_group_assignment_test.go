package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAppGroupAssignment_crud(t *testing.T) {
	resourceName0 := fmt.Sprintf("%s.test.0", resources.OktaIDaaSAppGroupAssignment)
	resourceName1 := fmt.Sprintf("%s.test.1", resources.OktaIDaaSAppGroupAssignment)
	resourceName3 := fmt.Sprintf("%s.test3", resources.OktaIDaaSAppGroupAssignment)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppGroupAssignment, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	// TF concurrency will cause this test to flap if the groups and assigned
	// priorities aren't executed in proper order
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "priority"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "priority"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName3),
					resource.TestCheckResourceAttrSet(resourceName3, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName3, "group_id"),
					resource.TestCheckResourceAttr(resourceName3, "priority", "3"),
					resource.TestCheckResourceAttr(resourceName3, "profile", "{}"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "priority"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "priority"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName3),
					resource.TestCheckResourceAttrSet(resourceName3, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName3, "group_id"),
					resource.TestCheckResourceAttr(resourceName3, "priority", "4"),
					resource.TestCheckResourceAttr(resourceName3, "profile", "{}"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName0),
					resource.TestCheckResourceAttrSet(resourceName0, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName0, "priority"),
					resource.TestCheckResourceAttr(resourceName0, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName1),
					resource.TestCheckResourceAttrSet(resourceName1, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "priority"),
					resource.TestCheckResourceAttr(resourceName1, "profile", "{}"),
					ensureAppGroupAssignmentExists(resourceName3),
					resource.TestCheckResourceAttrSet(resourceName3, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName3, "group_id"),
					resource.TestCheckResourceAttr(resourceName3, "priority", "3"),
					resource.TestCheckResourceAttr(resourceName3, "profile", "{}"),
				),
			},
			{
				ResourceName:      resourceName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName3]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName3)
					}
					appID := rs.Primary.Attributes["app_id"]
					groupID := rs.Primary.Attributes["group_id"]
					return fmt.Sprintf("%s/%s", appID, groupID), nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import schema into state")
					}
					return nil
				},
			},
		},
	})
}

func TestAccResourceOktaAppGroupAssignment_retain(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppGroupAssignment)
	appName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)
	groupName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroup)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppGroupAssignment, t.Name())
	retainAssignment := mgr.GetFixtures("retain_assignment.tf", t)
	retainAssignmentDestroy := mgr.GetFixtures("retain_assignment_destroy.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: retainAssignment,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttr(resourceName, "retain_assignment", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile", "{}"),
				),
			},
			{
				Config: retainAssignmentDestroy,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceNotExists(resourceName),
					ensureAppGroupAssignmentRetained(appName, groupName),
				),
			},
		},
	})
}

func TestAccResourceOktaAppGroupAssignment_timeouts(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppGroupAssignment, t.Name())
	resourceName0 := fmt.Sprintf("%s.test.0", resources.OktaIDaaSAppGroupAssignment)
	config := `
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
}

resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

locals {
  group_ids = tolist([okta_group.test.id])
}

resource "okta_app_group_assignment" "test" {
  count = length(local.group_ids)

  app_id   = okta_app_oauth.test.id
  group_id = local.group_ids[count.index]
  priority = count.index

  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}
`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName0, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName0, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName0, "timeouts.update", "30m"),
				),
			},
		},
	})
}

func ensureAppGroupAssignmentExists(resourceName string) resource.TestCheckFunc {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return func(s *terraform.State) error {
			return nil
		}
	}
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", resourceName)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return missingErr
		}

		appID := rs.Primary.Attributes["app_id"]
		groupID := rs.Primary.Attributes["group_id"]
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()

		g, _, err := client.Application.GetApplicationGroupAssignment(context.Background(), appID, groupID, nil)
		if err != nil {
			return err
		} else if g == nil {
			return missingErr
		}

		return nil
	}
}

func ensureAppGroupAssignmentRetained(appName, groupName string) resource.TestCheckFunc {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return func(s *terraform.State) error {
			return nil
		}
	}
	return func(s *terraform.State) error {
		notFound := "resource not found: %s"
		// app group assignment has been removed from state, so use app and group to query okta
		appRes, ok := s.RootModule().Resources[appName]
		if !ok {
			return fmt.Errorf(notFound, appName)
		}

		groupRes, ok := s.RootModule().Resources[groupName]
		if !ok {
			return fmt.Errorf(notFound, groupName)
		}

		appID := appRes.Primary.ID
		groupID := groupRes.Primary.ID
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()

		g, _, err := client.Application.GetApplicationGroupAssignment(context.Background(), appID, groupID, nil)
		if err != nil {
			return err
		} else if g == nil {
			return fmt.Errorf("application group assignment not found for app ID, group ID: %s, %s", appID, groupID)
		}
		return nil
	}
}
