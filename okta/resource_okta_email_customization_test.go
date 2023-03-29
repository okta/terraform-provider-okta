package okta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaEmailCustomization_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.forgot_password_en", emailCustomization)
	mgr := newFixtureManager(emailCustomization)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	updatedConfigChangeIsDefault := mgr.GetFixtures("updated_change_is_default.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceEmailCustomizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "id"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "brand_id"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "language", "en"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "is_default", "true"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "subject", "Forgot Password"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "body", "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "links"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "language", "es"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "is_default", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "language", "en"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "is_default", "true"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "body", "Hello $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "language", "es"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "is_default", "false"),
				),
			},
			{
				Config: updatedConfigChangeIsDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "language", "en"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "is_default", "false"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "language", "es"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "is_default", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName)
					}
					ID := rs.Primary.Attributes["id"]
					brandID := rs.Primary.Attributes["brand_id"]
					templateName := rs.Primary.Attributes["template_name"]
					return fmt.Sprintf("%s/%s/%s", ID, brandID, templateName), nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import schema into state")
					}
					return nil
				},
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

		_, resp, err := getOktaV3ClientFromMetadata(testAccProvider.Meta()).CustomizationApi.GetEmailCustomization(context.Background(), brandID, templateName, ID).Execute()
		if err != nil || resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("email customization still exists, ID %q, brandID %q, templateName: %q", ID, brandID, templateName)
	}
	return nil
}
