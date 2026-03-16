package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		CheckDestroy:             nil,
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
