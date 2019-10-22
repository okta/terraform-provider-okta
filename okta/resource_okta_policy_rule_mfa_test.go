package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func deleteMfaPolicyRules(client *testClient) error {
	return deletePolicyRulesByType(mfaPolicyType, client)
}

func TestAccOktaMfaPolicyRule_crud(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaMfaPolicyRule(ri)
	updatedConfig := testOktaMfaPolicyRuleUpdated(ri)
	resourceName := buildResourceFQN(policyRuleMfa, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(policyRuleMfa),
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
				),
			},
		},
	})
}

func testOktaMfaPolicyRule(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "MFA_ENROLL"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, policyRuleMfa, name, rInt, name)
}

func testOktaMfaPolicyRuleUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "MFA_ENROLL"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "INACTIVE"
	enroll	 = "LOGIN"
}
`, rInt, policyRuleMfa, name, rInt, name)
}
