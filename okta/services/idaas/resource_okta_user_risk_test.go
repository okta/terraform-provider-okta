package idaas_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaUserRisk_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUserRisk)
	mgr := newFixtureManager("resources", resources.OktaIDaaSUserRisk, t.Name())

	config := `
	resource "okta_user" "test" {
	first_name = "TestAcc"
	last_name  = "Smith"
	login      = "testAcc-replace_with_uuid@example.com"
	email      = "testAcc-replace_with_uuid@example.com"
	}

	resource "okta_user_risk" "test" {
	user_id    = okta_user.test.id
	risk_level = "LOW"
	}
	`

	updatedConfig := `
	resource "okta_user" "test" {
	first_name = "TestAcc"
	last_name  = "Smith"
	login      = "testAcc-replace_with_uuid@example.com"
	email      = "testAcc-replace_with_uuid@example.com"
	}

	resource "okta_user_risk" "test" {
	user_id    = okta_user.test.id
	risk_level = "HIGH"
	}
	`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserRiskTestUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "LOW"),
				),
			},
			{
				Config: mgr.ConfigReplace(updatedConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "HIGH"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataSourceOktaUserRisk_read(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUserRisk)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSUserRisk)
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSUserRisk, t.Name())

	config := `
	resource "okta_user" "test" {
	first_name = "TestAcc"
	last_name  = "Smith"
	login      = "testAcc-replace_with_uuid@example.com"
	email      = "testAcc-replace_with_uuid@example.com"
	}

	resource "okta_user_risk" "test" {
	user_id    = okta_user.test.id
	risk_level = "HIGH"
	}

	data "okta_user_risk" "test" {
	user_id = okta_user_risk.test.user_id
	}
	`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserRiskTestUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "HIGH"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_id", resourceName, "user_id"),
					resource.TestCheckResourceAttr(dataSourceName, "risk_level", "HIGH"),
				),
			},
		},
	})
}

func checkUserRiskTestUserDestroy(s *terraform.State) error {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	for _, r := range s.RootModule().Resources {
		if r.Type != "okta_user" {
			continue
		}
		if _, resp, err := client.User.GetUser(context.Background(), r.Primary.ID); err != nil {
			if resp != nil && resp.Response.StatusCode == http.StatusNotFound {
				continue
			}
			return fmt.Errorf("[ERROR] Error Getting User in Okta: %v", err)
		}
		return fmt.Errorf("user still exists")
	}
	return nil
}
