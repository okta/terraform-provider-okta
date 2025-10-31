package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// TestAccResourceOktaEmailCustomization_crud demonstrates having a default `en`
// customization and dependant `es` customization in step 1. Step 2 demonstrates
// flipping those default dependencies to have `es` be the default. Step 3 is
// flipping back to the state of step 1. This ACC tests full CRUD operations on
// the okta_email_customization resource and managing `is_default` dependencies
// with the `depends_on` meta argument.
func TestAccResourceOktaEmailCustomization_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.forgot_password_en", resources.OktaIDaaSEmailCustomization)
	mgr := newFixtureManager("resources", resources.OktaIDaaSEmailCustomization, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceEmailCustomizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "id"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "brand_id"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "language", "en"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "is_default", "true"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "subject", "Account password reset"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "body", "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "links"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "language", "es"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "is_default", "false"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "subject", "Restablecimiento de contraseña de cuenta"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "body", "Hola $$user.firstName,<br/><br/>Haga clic en este enlace para restablecer tu contraseña: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_es", "links"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "id"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "brand_id"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "language", "en"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "is_default", "false"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "subject", "Account password reset"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "body", "Hello $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "links"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "language", "es"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "is_default", "true"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "subject", "Restablecimiento de contraseña de cuenta"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "body", "Qué tal $$user.firstName,<br/><br/>Haga clic en este enlace para restablecer tu contraseña: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_es", "links"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "id"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "brand_id"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "language", "en"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "is_default", "true"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "subject", "Account password reset"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_en", "body", "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_en", "links"),

					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "template_name", "ForgotPassword"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "language", "es"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "is_default", "false"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "subject", "Restablecimiento de contraseña de cuenta"),
					resource.TestCheckResourceAttr("okta_email_customization.forgot_password_es", "body", "Hola $$user.firstName,<br/><br/>Haga clic en este enlace para restablecer tu contraseña: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("okta_email_customization.forgot_password_es", "links"),
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

func checkResourceEmailCustomizationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resources.OktaIDaaSEmailCustomization {
			continue
		}
		ID := rs.Primary.ID
		brandID := rs.Primary.Attributes["brand_id"]
		templateName := rs.Primary.Attributes["template_name"]
		t := &testing.T{}
		testClient := IDaaSClientForTest(t)
		client := testClient.OktaSDKClientV3()

		ctx := context.Background()

		_, _, err := client.CustomizationAPI.ListEmailCustomizations(ctx, brandID, templateName).Execute()
		if err != nil {
			return fmt.Errorf("failed to delete email customization ID %q, brandID %q, templateName: %q", ID, brandID, templateName)
		}
		_, err = client.CustomizationAPI.DeleteEmailCustomization(ctx, brandID, templateName, ID).Execute()
		if err != nil {
			return fmt.Errorf("failed to delete email customization ID %q, brandID %q, templateName: %q", ID, brandID, templateName)
		}

		_, resp, _ := client.CustomizationAPI.GetEmailCustomization(ctx, brandID, templateName, ID).Execute()
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("email customization still exists, ID %q, brandID %q, templateName: %q", ID, brandID, templateName)
	}
	return nil
}
