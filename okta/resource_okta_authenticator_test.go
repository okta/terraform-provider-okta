package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaAuthenticator_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", authenticator)
	mgr := newFixtureManager(authenticator)
	config := mgr.GetFixtures("security_question.tf", ri, t)
	configUpdated := mgr.GetFixtures("security_question_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
					testAuthenticatorSettings(resourceName, `{"allowedFor" : "recovery"}`),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
					testAuthenticatorSettings(resourceName, `{"allowedFor" : "any"}`),
				),
			},
			{
				Config: config,
			},
		},
	})
}

func testAuthenticatorSettings(name, expectedSettingsJSON string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		actualSettingsJSON := rs.Primary.Attributes["settings"]
		eq := areJSONStringsEqual(expectedSettingsJSON, actualSettingsJSON)
		if !eq {
			return fmt.Errorf("attribute 'settings' expected %q, got %q", expectedSettingsJSON, actualSettingsJSON)
		}
		return nil
	}
}
