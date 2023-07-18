package okta

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaOrgSupport(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", orgSupport)
	mgr := newFixtureManager(orgSupport, t.Name())
	config := mgr.GetFixtures("standard.tf", t)
	updatedConfig := mgr.GetFixtures("extended.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkSupportDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "ENABLED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "extend_by", "1"),
				),
			},
		},
	})
}

func checkSupportDestroy(s *terraform.State) error {
	if isVCRPlayMode() {
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != orgSupport {
			continue
		}
		client := oktaClientForTest()
		support, _, err := client.OrgSetting.GetOrgOktaSupportSettings(context.Background())
		if err != nil {
			return err
		}
		if support.Support == "ENABLED" {
			return errors.New("support is still enabled")
		}
	}
	return nil
}
