package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

func TestAccOktaSecurePasswordStoreApplication(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSecurePasswordStoreConfig(ri)
	updatedConfig := buildTestSecurePasswordStoreConfigUpdated(ri)
	resourceName := buildResourceFQN(securePasswordStoreApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(securePasswordStoreApp, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "url", "http://test.com"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "pass"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "EDIT_USERNAME_AND_PASSWORD"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://test.com/changed"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "user1"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "pass1"),
					resource.TestCheckResourceAttr(resourceName, "reveal_password", "false"),
					resource.TestCheckResourceAttr(resourceName, "shared_password", "sharedpassword"),
					resource.TestCheckResourceAttr(resourceName, "shared_username", "sharedusername"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "SHARED_USERNAME_AND_PASSWORD"),
				),
			},
		},
	})
}
func TestAccOktaSecurePasswordStoreApplicationCredsSchemes(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSecurePasswordStoreConfigAdmin(ri)
	updatedConfig := buildTestSecurePasswordStoreConfigExternalSync(ri)
	resourceName := buildResourceFQN(securePasswordStoreApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(securePasswordStoreApp, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
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

func buildTestSecurePasswordStoreConfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label         	 	= "%s"
  username_field		= "user"
  password_field		= "pass"
  url					= "http://test.com"
}
`, securePasswordStoreApp, name, name)
}

func buildTestSecurePasswordStoreConfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       			= "%s"
  status 	  			= "INACTIVE"
  username_field		= "user1"
  password_field		= "pass1"
  url					= "http://test.com/changed"
  reveal_password		= false
  credentials_scheme 	= "SHARED_USERNAME_AND_PASSWORD"
  shared_username 		= "sharedusername"
  shared_password		= "sharedpassword"
}
`, securePasswordStoreApp, name, name)
}

func buildTestSecurePasswordStoreConfigAdmin(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       			= "%s"
  username_field		= "user"
  password_field		= "pass"
  url					= "http://test.com"
  credentials_scheme 	= "ADMIN_SETS_CREDENTIALS"
}
`, securePasswordStoreApp, name, name)
}

func buildTestSecurePasswordStoreConfigExternalSync(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  status 	  			= "INACTIVE"
  label       			= "%s"
  username_field		= "user"
  password_field		= "pass"
  url					= "http://test.com"
  credentials_scheme 	= "EXTERNAL_PASSWORD_SYNC"
}
`, securePasswordStoreApp, name, name)
}
