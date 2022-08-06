package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppThreeFieldApplication_crud(t *testing.T) {
	mgr := newFixtureManager(appThreeField, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	updatedCreds := mgr.GetFixtures("updated_credentials.tf", t)
	resourceName := fmt.Sprintf("%s.test", appThreeField)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appThreeField, createDoesAppExist(okta.NewSwaThreeFieldApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaThreeFieldApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "button_selector", "btn"),
					resource.TestCheckResourceAttr(resourceName, "username_selector", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_selector", "pass"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_selector", "third"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_value", "third"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaThreeFieldApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "button_selector", "btn1"),
					resource.TestCheckResourceAttr(resourceName, "username_selector", "user1"),
					resource.TestCheckResourceAttr(resourceName, "password_selector", "pass1"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://example.com"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_selector", "mfa"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_value", "mfa"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedCreds,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaThreeFieldApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "button_selector", "btn1"),
					resource.TestCheckResourceAttr(resourceName, "username_selector", "user1"),
					resource.TestCheckResourceAttr(resourceName, "password_selector", "pass1"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://example.com"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_selector", "mfa"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_value", "mfa"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "SHARED_USERNAME_AND_PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "shared_username", buildResourceName(mgr.Seed)+"@example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_password"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}
