package okta

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThreatInsightSettings(t *testing.T) {
	mgr := newFixtureManager(threatInsightSettings, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", threatInsightSettings)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaThreatInsightSettingsDestroy(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "none"),
					resource.TestCheckResourceAttr(resourceName, "network_excludes.#", "0"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "block"),
					resource.TestCheckResourceAttr(resourceName, "network_excludes.#", "1"),
				),
			},
		},
	})
}

func checkOktaThreatInsightSettingsDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != threatInsightSettings {
				continue
			}
			conf, _, err := getOktaClientFromMetadata(testAccProvider.Meta()).ThreatInsightConfiguration.GetCurrentConfiguration(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get threat insight configuration: %v", err)
			}
			if len(conf.ExcludeZones) > 0 {
				return errors.New("excluded zones list should be empty")
			}
			if conf.Action != "none" {
				return errors.New("action should be 'none'")
			}
		}
		return nil
	}
}
