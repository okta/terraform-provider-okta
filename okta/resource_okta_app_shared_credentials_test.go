package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppSharedCredentials_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSharedCredentials)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appSharedCredentials)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appSharedCredentials, createDoesAppExist(okta.NewBrowserPluginApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewBrowserPluginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login-updated.html"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://example.com/redirect_url"),
					resource.TestCheckResourceAttr(resourceName, "checkbox", "checkbox_red"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template", "user.firstName"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template_suffix", "hello"),
					resource.TestCheckResourceAttr(resourceName, "shared_password", "sharedpass"),
					resource.TestCheckResourceAttr(resourceName, "shared_username", "sharedusername"),
					resource.TestCheckResourceAttr(resourceName, "accessibility_self_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "accessibility_error_redirect_url", "https://example.com/redirect_url_1"),
					resource.TestCheckResourceAttr(resourceName, "accessibility_login_redirect_url", "https://example.com/redirect_url_2"),
					resource.TestCheckResourceAttr(resourceName, "auto_submit_toolbar", "true"),
					resource.TestCheckResourceAttr(resourceName, "hide_ios", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewBrowserPluginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login-updated"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username-updated"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password-updated"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login-updated.html"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://example.com/redirect_url"),
					resource.TestCheckResourceAttr(resourceName, "checkbox", "checkbox_red-updated"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template", "user.firstName"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template_suffix", "moas"),
					resource.TestCheckResourceAttr(resourceName, "shared_password", "sharedpass22"),
					resource.TestCheckResourceAttr(resourceName, "shared_username", "sharedusername22"),
					resource.TestCheckResourceAttr(resourceName, "accessibility_self_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "accessibility_error_redirect_url", "https://example.com/redirect_url_1"),
					resource.TestCheckResourceAttr(resourceName, "accessibility_login_redirect_url", "https://example.com/redirect_url_2"),
					resource.TestCheckResourceAttr(resourceName, "auto_submit_toolbar", "true"),
					resource.TestCheckResourceAttr(resourceName, "hide_ios", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}
