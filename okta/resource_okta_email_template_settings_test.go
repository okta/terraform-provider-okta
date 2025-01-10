package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaEmailTemplateSettings(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", emailTemplateSettings)
	mgr := newFixtureManager("resources", emailTemplateSettings, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "recipients", "NO_USERS"),
					resource.TestCheckResourceAttr(resourceName, "template_name", "UserActivation"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "recipients", "ADMINS_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "template_name", "UserActivation"),
				),
			},
		},
	})
}
