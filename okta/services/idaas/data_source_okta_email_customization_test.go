package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEmailCustomization_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSEmailCustomization, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
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
					resource.TestCheckResourceAttr("data.okta_email_customization.forgot_password_en", "is_default", "true"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "subject"),
					resource.TestCheckResourceAttr("data.okta_email_customization.forgot_password_en", "subject", "Forgot Password"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "body"),
					resource.TestCheckResourceAttr("data.okta_email_customization.forgot_password_en", "body", "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"),
					resource.TestCheckResourceAttrSet("data.okta_email_customization.forgot_password_en", "links"),
				),
			},
		},
	})
}
