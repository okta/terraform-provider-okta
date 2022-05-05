package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaEmailCustomization_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(emailCustomization)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "id"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "brand_id"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "template_name"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "language"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "language", "cs"), // setting the language to Czech for testing
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "is_default"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "is_default", "false"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "subject"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "subject", "Forgot Password"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "body"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "body", "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "links"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "id"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "brand_id"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "template_name"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "language"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "language", "cs"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "is_default"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "is_default", "false"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "subject"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "subject", "Forgot Password"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "body"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en_alt", "body", "Hello $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en_alt", "links"),
				),
			},
		},
	})
}
