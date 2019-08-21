package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func deletePolicyRuleIdpDiscovery(client *testClient) error {
	return deletePolicyRulesByType(idpDiscovery, client)
}

func TestAccOktaPolicyRuleIdpDiscovery(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyRuleIdpDiscovery)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_domain.tf", ri, t)
	deactivatedConfig := mgr.GetFixtures("basic_deactivated.tf", ri, t)
	appIncludeConfig := mgr.GetFixtures("app_include.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyRuleIdpDiscovery)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(policyRuleIdpDiscovery),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "2"),
				),
			},
			{
				Config: deactivatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "2"),
				),
			},
			{
				Config: appIncludeConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "app_include.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "idp_type", "SAML2"),
				),
			},
		},
	})
}
