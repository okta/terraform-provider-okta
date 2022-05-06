package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaEmailCustomization_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(emailCustomization)
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "customization_id"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "id"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "brand_id"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "template_name"),
					resource.TestCheckResourceAttr("data.okta_email_customization.forgot_password_en", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "language"),
					resource.TestCheckResourceAttr("data.okta_email_customization.forgot_password_en", "language", "en"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "is_default"),
					resource.TestCheckResourceAttr("data.okta_email_customization.forgot_password_en", "is_default", "false"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "subject"),
					resource.TestCheckResourceAttr("data.okta_email_customization.forgot_password_en", "subject", "Stuff"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "body"),
					resource.TestCheckResourceAttr("data.okta_email_customization.forgot_password_en", "body", "Hi $$user.firstName,<br/><br/>Blah blah $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "links"),
				),
			},
		},
	})
}
