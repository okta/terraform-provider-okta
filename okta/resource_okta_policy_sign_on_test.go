package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/sdk"
)

func deleteSignOnPolicies(client *testClient) error {
	return deletePolicyByType(sdk.SignOnPolicyType, client)
}

func TestAccOktaPolicySignOn_defaultError(t *testing.T) {
	mgr := newFixtureManager(policySignOn, t.Name())
	config := testOktaPolicySignOnDefaultErrors(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createPolicyCheckDestroy(policySignOn),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Default Policy is immutable"),
			},
		},
	})
}

func TestAccOktaPolicySignOn_crud(t *testing.T) {
	mgr := newFixtureManager(policySignOn, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_inactive.tf", t)
	renamedConfig := mgr.GetFixtures("basic_renamed.tf", t)
	resourceName := fmt.Sprintf("%s.test", policySignOn)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createPolicyCheckDestroy(policySignOn),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy Updated"),
				),
			},
			{
				Config: renamedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAccUpdated_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
		},
	})
}

func testOktaPolicySignOnDefaultErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "Default Policy"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy"
}
`, policySignOn, name)
}
