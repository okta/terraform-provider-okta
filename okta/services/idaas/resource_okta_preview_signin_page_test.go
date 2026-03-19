package idaas_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccResourceOktaPreviewSignInPage_crud(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_preview_signin_page", t.Name())

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("basic.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_preview_signin_page.test", "page_content", "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\">\n<html>\n<head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n    <meta name=\"robots\" content=\"noindex,nofollow\" />\n    <!-- Styles generated from theme -->\n    <link href=\"{{themedStylesUrl}}\" rel=\"stylesheet\" type=\"text/css\">\n    <!-- Favicon from theme -->\n    <link rel=\"shortcut icon\" href=\"{{faviconUrl}}\" type=\"image/x-icon\"/>\n\n    <title>{{pageTitle}}</title>\n    {{{SignInWidgetResources}}}\n\n    <style nonce=\"{{nonceValue}}\">\n        #login-bg-image-id {\n            background-image: {{bgImageUrl}}\n        }\n    </style>\n</head>\n<body>\n    <div id=\"login-bg-image-id\" class=\"login-bg-image tb--background\"></div>\n    <div id=\"okta-login-container\"></div>\n\n    <!--\n        \"OktaUtil\" defines a global OktaUtil object\n        that contains methods used to complete the Okta login flow.\n     -->\n    {{{OktaUtil}}}\n\n    <script type=\"text/javascript\" nonce=\"{{nonceValue}}\">\n        // \"config\" object contains default widget configuration\n        // with any custom overrides defined in your admin settings.\n        var config = OktaUtil.getSignInWidgetConfig();\n\n        // Render the Okta Sign-In Widget\n        var oktaSignIn = new OktaSignIn(config);\n        oktaSignIn.renderEl({ el: '#okta-login-container' },\n            OktaUtil.completeLogin,\n            function(error) {\n                // Logs errors that occur when configuring the widget.\n                // Remove or replace this with your own custom error handler.\n                console.log(error.message, error);\n            }\n        );\n    </script>\n</body>\n</html>\n"),
					resource.TestCheckResourceAttr("okta_preview_signin_page.test", "widget_version", "^6"),
					resource.TestCheckResourceAttr("okta_preview_signin_page.test", "widget_customizations.widget_generation", "G2"),
					resource.TestCheckResourceAttr("okta_preview_signin_page.test", "content_security_policy_setting.mode", "report_only"),
					resource.TestCheckResourceAttr("okta_preview_signin_page.test", "content_security_policy_setting.report_uri", ""),
					resource.TestCheckResourceAttr("okta_preview_signin_page.test", "content_security_policy_setting.src_list.#", "2"),
				),
			},
			// Regression test for https://github.com/okta/terraform-provider-okta/issues/2201
			// Importing the resource must preserve brand_id so that subsequent reads succeed.
			{
				ResourceName:      "okta_preview_signin_page.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["okta_preview_signin_page.test"]
					if !ok {
						return "", fmt.Errorf("failed to find okta_preview_signin_page.test")
					}
					return rs.Primary.Attributes["brand_id"], nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import resource into state")
					}
					if s[0].Attributes["brand_id"] == "" {
						return errors.New("brand_id is empty after import; import state bug not fixed")
					}
					if s[0].ID == "" {
						return errors.New("id is empty after import; import state bug not fixed")
					}
					return nil
				},
			},
		},
	})
}
