package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// TestAccDataSourceOktaPolicyRulePassword_read verifies that the data source can look up a
// password policy rule by name and correctly exposes all SSPR requirement fields, including
// method_constraints, primary_methods, step_up_enabled, and step_up_methods.
func TestAccDataSourceOktaPolicyRulePassword_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSPolicyRulePassword, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := "data.okta_policy_rule_password." + acctest.BuildResourceName(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "password_change", "ALLOW"),
					resource.TestCheckResourceAttr(dataSourceName, "password_reset", "ALLOW"),
					resource.TestCheckResourceAttr(dataSourceName, "password_unlock", "ALLOW"),
					resource.TestCheckResourceAttr(dataSourceName, "password_reset_access_control", "LEGACY"),
					// People conditions
					resource.TestCheckResourceAttr(dataSourceName, "users_included.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "users_excluded.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "groups_included.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "groups_excluded.#", "1"),
					// SSPR requirement fields
					resource.TestCheckResourceAttr(dataSourceName, "password_reset_requirement.0.step_up_enabled", "true"),
					resource.TestCheckTypeSetElemAttr(dataSourceName, "password_reset_requirement.0.step_up_methods.*", "security_question"),
					resource.TestCheckTypeSetElemAttr(dataSourceName, "password_reset_requirement.0.primary_methods.*", "otp"),
					resource.TestCheckTypeSetElemAttr(dataSourceName, "password_reset_requirement.0.primary_methods.*", "email"),
					resource.TestCheckResourceAttr(dataSourceName, "password_reset_requirement.0.method_constraints.0.method", "otp"),
				),
			},
		},
	})
}
