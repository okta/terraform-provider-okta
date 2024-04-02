package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaPolicyRuleIdpDiscovery_crud(t *testing.T) {
	mgr := newFixtureManager("resources", policyRuleIdpDiscovery, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_domain.tf", t)
	deactivatedConfig := mgr.GetFixtures("basic_deactivated.tf", t)

	mgr2 := newFixtureManager("resources", policyRuleIdpDiscovery, t.Name())
	appIncludeConfig := mgr2.GetFixtures("app_include.tf", t)
	appExcludeConfig := mgr2.GetFixtures("app_exclude_platform.tf", t)
	resourceName := fmt.Sprintf("%s.test", policyRuleIdpDiscovery)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRuleIdpDiscovery),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_type", "ATTRIBUTE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_type", "IDENTIFIER"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "2"),
				),
			},
			{
				Config: deactivatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_type", "IDENTIFIER"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "2"),
				),
			},
			{
				Config: appIncludeConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "app_include.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "idp_type", "OKTA"),
				),
			},
			{
				Config: appExcludeConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "app_exclude.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "idp_type", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceOktaPolicyRuleIdpDiscoveryDynamic(t *testing.T) {
	mgr := newFixtureManager("resources", policyRuleIdpDiscovery, t.Name())
	config := mgr.GetFixtures("dynamic.tf", t)
	resourceName := fmt.Sprintf("%s.test", policyRuleIdpDiscovery)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRuleIdpDiscovery),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "provider_expression", "login.identifier.substringAfter('@')"),
					resource.TestCheckResourceAttr(resourceName, "selection_type", "DYNAMIC"),
				),
			},
		},
	})
}
