package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaEmailCustomizations_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", emailCustomizations, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.#"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.#", "2"),

					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.0.id"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.0.language"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.0.language", "en"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.0.is_default"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.0.is_default", "true"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.0.subject"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.0.subject", "Forgot Password"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.0.body"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.0.body", "Hi $$user.firstName,<br/><br/>Blah blah $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.0.links"),

					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.1.id"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.1.language"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.1.language", "es"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.1.is_default"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.1.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.1.subject"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.1.subject", "Has olvidado tu contraseña"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.1.body"),
					resource.TestCheckResourceAttr("data.okta_email_customizations.forgot_password", "email_customizations.1.body", "Hola $$user.firstName,<br/><br/>Puedo ir al bano $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("data.okta_email_customizations.forgot_password", "email_customizations.1.links"),
				),
			},
		},
	})
}
