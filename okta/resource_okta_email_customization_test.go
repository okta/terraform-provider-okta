package okta

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaEmailCustomization_crud(t *testing.T) {
	mgr := newFixtureManager(emailCustomization, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceEmailCustomizationDestroy,
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

func createCheckResourceEmailCustomizationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != emailCustomization {
			continue
		}
		ID := rs.Primary.ID
		brandID := rs.Primary.Attributes["brand_id"]
		templateName := rs.Primary.Attributes["template_name"]

		_, resp, err := getOktaClientFromMetadata(testAccProvider.Meta()).Brand.GetEmailTemplateCustomization(context.Background(), brandID, templateName, ID)
		if err != nil || resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("email customization still exists, ID %q, brandID %q, templateName: %q", ID, brandID, templateName)
	}
	return nil
}
