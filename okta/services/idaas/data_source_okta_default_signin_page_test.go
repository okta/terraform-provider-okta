package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaDefaultSigninPage_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_default_signin_page", t.Name())
	resourceName := fmt.Sprintf("data.%s.test", "okta_default_signin_page")

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("datasource.tf", t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "widget_version", "^7"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.sign_in_label", "Sign In"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.username_label", "Username"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.password_label", "Password"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.show_password_visibility_toggle", "true"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.show_user_identifier", "true"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.forgot_password_label", "Forgot password?"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.unlock_account_label", "Unlock account?"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.help_label", "Help"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.classic_recovery_flow_email_or_username_label", "Email or Username"),
					resource.TestCheckResourceAttr(resourceName, "widget_customizations.widget_generation", "G2"),
				),
			},
		},
	})
}
