package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaPolicyRuleProfileEnrollment(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyRuleProfileEnrollment)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyRuleProfileEnrollment)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createRuleCheckDestroy(policyRuleProfileEnrollment),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "unknown_user_action", "REGISTER"),
					resource.TestCheckResourceAttr(resourceName, "email_verification", "true"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.#", "0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "unknown_user_action", "REGISTER"),
					resource.TestCheckResourceAttr(resourceName, "email_verification", "true"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttrSet(resourceName, "inline_hook_id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_group_id"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.0.name", "email"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.1.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.2.name", "t-shirt"),
				),
			},
		},
	})
}
