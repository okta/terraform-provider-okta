package okta

import (
	"fmt"
	"regexp"
	"testing"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

func deleteSignOnPolicyRules(artClient *articulateOkta.Client, client *okta.Client) error {
	return deletePolicyRulesByType(signOnPolicyType, artClient, client)
}

func TestAccOktaPolicyRuleDefaultErrors(t *testing.T) {
	config := testOktaPolicyRuleSignOnDefaultErrors(acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("You cannot edit a default Policy Rule"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaPolicyRulesRename(t *testing.T) {
	ri := acctest.RandInt()
	updatedName := fmt.Sprintf("%s-changed-%d", testResourcePrefix, ri)
	config := testOktaPolicyRuleSignOn(ri)
	updatedConfig := testOktaPolicyRuleSignOnRename(updatedName, ri)
	resourceName := buildResourceFQN(signOnPolicyRule, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(signOnPolicyRule),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
		},
	})
}

func TestAccOktaPolicyRulesNewPolicy(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOn(ri)
	updatedConfig := testOktaPolicyRuleSignOnNewPolicy(ri)
	resourceName := buildResourceFQN(signOnPolicyRule, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(signOnPolicyRule),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
				),
			},
		},
	})
}
func TestAccOktaPolicyRuleSignOn(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOn(ri)
	updatedConfig := testOktaPolicyRuleSignOnUpdated(ri)
	resourceName := buildResourceFQN(signOnPolicyRule, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(signOnPolicyRule),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "access", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "session_idle", "240"),
					resource.TestCheckResourceAttr(resourceName, "session_lifetime", "240"),
					resource.TestCheckResourceAttr(resourceName, "session_persistent", "false"),
				),
			},
		},
	})
}

func TestAccOktaPolicyRuleSignOnPassErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOnPassErrors(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(signOnPolicyRule),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("config is invalid: .*: : invalid or unknown key: password_change"),
				PlanOnly:    true,
			},
		},
	})
}

func testOktaPolicyRuleSignOn(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, signOnPolicyType, signOnPolicyRule, name, rInt, name)
}

func testOktaPolicyRuleSignOnUpdated(rInt int) string {
	name := buildResourceName(rInt)

	// Adding a second resource here to test the priority preference
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "INACTIVE"
	access           = "DENY"
	session_idle      = 240
	session_lifetime  = 240
	session_persistent = false
}
`, rInt, signOnPolicyType, signOnPolicyRule, name, rInt, name)
}

func testOktaPolicyRuleSignOnDefaultErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
	policyid = "garbageID"
	name     = "Default Rule"
	status   = "ACTIVE"
}
`, signOnPolicyRule, name)
}

func testOktaPolicyRuleSignOnRename(updatedName string, rInt int) string {
	name := buildResourceName(rInt)
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, signOnPolicyType, signOnPolicyRule, name, rInt, updatedName)
}

func testOktaPolicyRuleSignOnNewPolicy(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	name        = "%s"
	status      = "ACTIVE"
	description = "Terraform Acceptance Test SignOn Policy"
}

resource "%s" "%s" {
  	policyid = "${okta_signon_policy.%s.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, signOnPolicyType, signOnPolicy, name, name, signOnPolicyRule, name, name, name)
}

func testOktaPolicyRuleSignOnPassErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  name     = "%s"
  status   = "ACTIVE"
  password_change = "DENY"
}
`, rInt, signOnPolicyType, signOnPolicyRule, name, rInt, name)
}
