package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

func TestAccAppSecurePasswordStoreApplication_credsSchemes(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSecurePasswordStore)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appSecurePasswordStore)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSecurePasswordStore, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "url", "http://test.com"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "pass"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "ADMIN_SETS_CREDENTIALS"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://test.com"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "pass"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "EXTERNAL_PASSWORD_SYNC"),
				),
			},
		},
	})
}
