package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func deleteMfaPolicyRules(client *testClient) error {
	return deletePolicyRulesByType(mfaPolicyType, client)
}

func TestAccOktaMfaPolicyRule(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaMfaPolicyRule(ri)
	updatedConfig := testOktaMfaPolicyRuleUpdated(ri)
	resourceName := buildResourceFQN(mfaPolicyRule, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(mfaPolicyRule),
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
data "okta_default_policies" "default-%d" {
	type = "MFA_ENROLL"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, mfaPolicyRule, name, rInt, name)
}

func testOktaMfaPolicyRuleUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "MFA_ENROLL"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "INACTIVE"
	enroll	 = "LOGIN"
}
`, rInt, mfaPolicyRule, name, rInt, name)
}
