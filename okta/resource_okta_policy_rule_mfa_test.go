package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func deleteMfaPolicyRules(client *testClient) error {
	return deletePolicyRulesByType(sdk.MfaPolicyType, client)
}

func TestAccOktaMfaPolicyRule_crud(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaMfaPolicyRule(ri)
	updatedConfig := testOktaMfaPolicyRuleUpdated(ri)
	resourceName := buildResourceFQN(policyRuleMfa, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createRuleCheckDestroy(policyRuleMfa),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
				),
			},
		},
	})
}

func testOktaMfaPolicyRule(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, sdk.MfaPolicyType, policyRuleMfa, name, rInt, name)
}

func testOktaMfaPolicyRuleUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "INACTIVE"
	enroll	 = "LOGIN"
}
`, rInt, sdk.MfaPolicyType, policyRuleMfa, name, rInt, name)
}
