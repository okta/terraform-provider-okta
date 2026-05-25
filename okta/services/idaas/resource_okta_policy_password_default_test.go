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
