package okta4ai_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaClientAuthSetting_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaClientAuthSetting, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaClientAuthSetting)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkClientAuthSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "OAUTH2"),
					resource.TestCheckResourceAttr(resourceName, "purposes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc-cred-%d", mgr.Seed)),
				),
			},
			{
				// Changing the configuration updates the resource in place.
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc-cred-updated-%d", mgr.Seed)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Write-only and request-only inputs: Read never returns them.
				ImportStateVerifyIgnore: []string{
					"resource_orn",
					"purposes",
					"oauth_client_configuration.client_credentials.client_secret",
				},
			},
		},
	})
}

func checkClientAuthSettingDestroy(s *terraform.State) error {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}
	return nil
}
