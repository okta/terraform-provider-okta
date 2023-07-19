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

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaThreatInsightSettingsDestroy,
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

// TestAccThreatInsightSettingsNetworkZoneOrdering https://github.com/okta/terraform-provider-okta/issues/1221
func TestAccThreatInsightSettingsNetworkZoneOrdering(t *testing.T) {
	mgr := newFixtureManager(threatInsightSettings, t.Name())
	resourceName := fmt.Sprintf("%s.test", threatInsightSettings)
	config := `
	resource "okta_network_zone" "a" {
		name     = "testAcc_replace_with_uuid-1"
		type     = "IP"
		gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
		proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
		status   = "ACTIVE"
	}
	resource "okta_network_zone" "b" {
		name     = "testAcc_replace_with_uuid-2"
		type     = "IP"
		gateways = ["2.2.3.4/24", "2.3.4.5-2.3.4.15"]
		proxies  = ["3.2.3.4/24", "3.3.4.5-3.3.4.15"]
		status   = "ACTIVE"
	}
	resource "okta_network_zone" "c" {
		name     = "testAcc_replace_with_uuid-3"
		type     = "IP"
		gateways = ["3.2.3.4/24", "2.3.4.5-2.3.4.15"]
		proxies  = ["4.2.3.4/24", "3.3.4.5-3.3.4.15"]
		status   = "ACTIVE"
	}
	resource "okta_threat_insight_settings" "test" {
		action           = "block"
		network_excludes = [okta_network_zone.a.id,okta_network_zone.b.id,okta_network_zone.c.id]
		#depends_on = [okta_network_zone.a, okta_network_zone.b, okta_network_zone.b]
	}
	`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaThreatInsightSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "block"),
					resource.TestCheckResourceAttr(resourceName, "network_excludes.#", "3"),
				),
			},
		},
	})
}

func checkOktaThreatInsightSettingsDestroy(s *terraform.State) error {
	if isVCRPlayMode() {
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != threatInsightSettings {
			continue
		}
		client := oktaClientForTest()
		conf, _, err := client.ThreatInsightConfiguration.GetCurrentConfiguration(context.Background())
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
