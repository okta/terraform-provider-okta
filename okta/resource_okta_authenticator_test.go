package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaAuthenticator_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", oktaAuthenticator)
	mgr := newFixtureManager("okta_authenticator")
	config := mgr.GetFixtures("okta_authenticator.tf", ri, t)
	activeConfig := mgr.GetFixtures("okta_authenticator_active.tf", ri, t)
	inactiveConfig := mgr.GetFixtures("okta_authenticator_inactive.tf", ri, t)
	// TODO settings tests
	// configSettingsAll5 := mgr.GetFixtures("okta_authenticator_settings_allowed_for_all_lifetime_5.tf", ri, t)
	// configSettingsRecovery10 := mgr.GetFixtures("okta_authenticator_settings_allowed_for_recovery_lifetime_10.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "key"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "settings"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
				),
			},
			{
				Config: activeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE")),
			},
			{
				Config: inactiveConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE")),
			},

			// TODO test updating settings when the okta-sdk-golang package adds
			// support for updating settings.
			// {
			// 	Config: configSettingsAll5,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceName, "settings", `{"allowedFor":"all","tokenLifetimeInMinutes":5}`)),
			// },
			// {
			// 	Config: configSettingsRecovery10,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceName, "settings", `{"allowedFor":"recovery","tokenLifetimeInMinutes":10}`)),
			// },
		},
	})
}
