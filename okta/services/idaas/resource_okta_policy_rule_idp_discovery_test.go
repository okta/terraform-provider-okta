package idaas_test

import (
	"fmt"
	"testing"

	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaPolicyRuleIdpDiscovery_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleIdpDiscovery, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_domain.tf", t)
	deactivatedConfig := mgr.GetFixtures("basic_deactivated.tf", t)
	appIncludeConfig := mgr.GetFixtures("app_include.tf", t)
	appExcludeConfig := mgr.GetFixtures("app_exclude_platform.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyRuleIdpDiscovery)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleIdpDiscovery),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_type", "ATTRIBUTE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_type", "IDENTIFIER"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "2"),
				),
			},
			{
				Config: deactivatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_type", "IDENTIFIER"),
					resource.TestCheckResourceAttr(resourceName, "user_identifier_patterns.#", "2"),
				),
			},
			{
				Config: appIncludeConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "app_include.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "idp_providers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "idp_providers.0.type", "OKTA"),
				),
			},
			{
				Config: appExcludeConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "app_exclude.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "idp_providers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "idp_providers.0.type", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.#", "1"),
				),
			},
		},
	})
}
