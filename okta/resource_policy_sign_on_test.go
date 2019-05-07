package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func deleteSignOnPolicies(client *testClient) error {
	return deletePolicyByType(signOnPolicyType, client)
}

func TestAccOktaPoliciesDefaultErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOnDefaultErrors(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createPolicyCheckDestroy(policySignOn),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Default Policy is immutable"),
			},
		},
	})
}

func TestAccOktaPolicySignOn(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policySignOn)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_inactive.tf", ri, t)
	renamedConfig := mgr.GetFixtures("basic_renamed.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policySignOn)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createPolicyCheckDestroy(policySignOn),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy Updated"),
				),
			},
			{
				Config: renamedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAccUpdated_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
		},
	})
}

func TestAccOktaPolicySignOnPassErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOnPassErrors(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createPolicyCheckDestroy(policySignOn),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("config is invalid: .* : invalid or unknown key: password_min_length"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaPolicySignOnAuthErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOnAuthErrors(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createPolicyCheckDestroy(policySignOn),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("config is invalid: .* : invalid or unknown key: auth_provider"),
				PlanOnly:    true,
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

func testOktaPolicySignOnPassErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
  password_min_length = 12
}
`, policySignOn, name, name)
}

func testOktaPolicySignOnAuthErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
  auth_provider = "ACTIVE_DIRECTORY"
}
`, policySignOn, name, name)
}
