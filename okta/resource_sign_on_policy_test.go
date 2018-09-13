package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/okta/okta-sdk-golang/okta"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func deleteSignOnPolicies(artClient *articulateOkta.Client, client *okta.Client) error {
	return deletePolicyByType(signOnPolicyType, artClient, client)
}

func TestAccOktaPolicies_defaultErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn_defaultErrors(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("You cannot edit a default Policy"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaPoliciesRename(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedName := fmt.Sprintf("%s-changed-%d", testResourcePrefix, ri)
	updatedConfig := testOktaPolicySignOnRename(updatedName, ri)
	resourceName := buildResourceFQN(signOnPolicy, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
		},
	})
}
func TestAccOktaPolicySignOn(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedConfig := testOktaPolicySignOn_updated(ri)
	resourceName := buildResourceFQN(signOnPolicy, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy Updated"),
				),
			},
		},
	})
}

func TestAccOktaPolicySignOn_passErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedConfig := testOktaPolicySignOn_passErrors(ri)
	resourceName := buildResourceFQN(signOnPolicy, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("config is invalid: .* : invalid or unknown key: password_min_length"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaPolicySignOn_authpErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedConfig := testOktaPolicySignOn_authpErrors(ri)
	resourceName := buildResourceFQN(signOnPolicy, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("config is invalid: .* : invalid or unknown key: auth_provider"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
		},
	})
}

func testOktaPolicySignOn(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy"
}
`, signOnPolicy, name, name)
}

func testOktaPolicySignOn_updated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_everyone_group" "everyone-%d" {}

resource "%s" "%s" {
  name        = "%s"
  status      = "INACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Updated"
  groups_included = [ "${data.okta_everyone_group.everyone-%d.id}" ]
}
`, rInt, signOnPolicy, name, name, rInt)
}

func testOktaPolicySignOn_defaultErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "Default Policy"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy"
}
`, signOnPolicy, name)
}

func testOktaPolicySignOnRename(updatedName string, rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
}
`, signOnPolicy, name, updatedName)
}

func testOktaPolicySignOn_passErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
  password_min_length = 12
}
`, signOnPolicy, name, name)
}

func testOktaPolicySignOn_authpErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
  auth_provider = "ACTIVE_DIRECTORY"
}
`, signOnPolicy, name, name)
}
