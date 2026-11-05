package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaDefaultPasswordPolicy_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyPasswordDefault, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyPasswordDefault)

	// NOTE needs the "Security Question" authenticator enabled on the org
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "sms_recovery", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "password_history_count", "5"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "sms_recovery", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "password_history_count", "0"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "sms_recovery", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "password_history_count", "0"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "sms_recovery", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "password_history_count", "5"),
				),
			},
		},
	})
}

// TestAccResourceOktaDefaultPasswordPolicy_issue_2804 reproduces GH-2804 / OKTA-1167884.
// When more than one password policy exists, the default policy's priority becomes 2.
// Earlier versions of the provider sent the read-only `priority` attribute back in the
// PUT body, which the Okta API rejects with E0000077. The default-policy update path
// must omit `priority` so this scenario succeeds.
func TestAccResourceOktaDefaultPasswordPolicy_issue_2804(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyPasswordDefault, t.Name())
	config := mgr.GetFixtures("issue_2804.tf", t)
	updatedConfig := mgr.GetFixtures("issue_2804_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyPasswordDefault)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "password_history_count", "5"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "password_history_count", "0"),
				),
			},
		},
	})
}

// TestAccResourceOktaDefaultPasswordPolicy_complexityExtras requires an OIE org
// with the maxConsecutiveRepeatingCharacters and oelStatement password
// complexity features enabled. The final step removes both attributes so the
// org's default policy is left without them.
func TestAccResourceOktaDefaultPasswordPolicy_complexityExtras(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyPasswordDefault, t.Name())
	config := mgr.GetFixtures("complexity_extras.tf", t)
	updatedConfig := mgr.GetFixtures("complexity_extras_updated.tf", t)
	removedConfig := mgr.GetFixtures("complexity_extras_removed.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyPasswordDefault)

	// The new attributes must not disturb the existing complexity settings,
	// which keep their schema defaults in every fixture.
	checkDefaultComplexity := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "password_min_length", "8"),
		resource.TestCheckResourceAttr(resourceName, "password_min_lowercase", "1"),
		resource.TestCheckResourceAttr(resourceName, "password_min_uppercase", "1"),
		resource.TestCheckResourceAttr(resourceName, "password_min_number", "1"),
		resource.TestCheckResourceAttr(resourceName, "password_min_symbol", "0"),
		resource.TestCheckResourceAttr(resourceName, "password_dictionary_lookup", "false"),
	)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "password_max_consecutive_repeating_characters", "3"),
					resource.TestCheckResourceAttr(resourceName, "password_oel_statement", "password.value.contains('okta')"),
					checkDefaultComplexity,
				),
			},
			{
				// Re-reading the policy must not produce a diff.
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_max_consecutive_repeating_characters", "5"),
					resource.TestCheckResourceAttr(resourceName, "password_oel_statement", "password.value.contains('okta') OR password.value.contains('terraform')"),
					checkDefaultComplexity,
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "password_max_consecutive_repeating_characters", "5"),
					resource.TestCheckResourceAttr(resourceName, "password_oel_statement", "password.value.contains('okta') OR password.value.contains('terraform')"),
				),
			},
			{
				// Removing both attributes from config clears them upstream.
				Config: removedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_max_consecutive_repeating_characters", "0"),
					resource.TestCheckResourceAttr(resourceName, "password_oel_statement", ""),
					checkDefaultComplexity,
				),
			},
			{
				Config:   removedConfig,
				PlanOnly: true,
			},
		},
	})
}
