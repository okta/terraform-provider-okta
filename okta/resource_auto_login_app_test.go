package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

// Test creation of a simple AWS SWA auto login app. The preconfigured apps are created by name.
func TestAccOktaAutoLoginApplicationPreconfig(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestAutoLoginConfigPreconfig(ri)
	updatedConfig := buildTestAutoLoginConfigPreconfigUpdated(ri)
	resourceName := buildResourceFQN(autoLoginApp, ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(autoLoginApp, createDoesAppExist(okta.NewAutoLoginApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
		},
	})
}

func TestAccOktaAutoLoginApplication(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestAutoLoginConfig(ri)
	updatedConfig := buildTestAutoLoginConfigUpdated(ri)
	resourceName := buildResourceFQN(autoLoginApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(autoLoginApp, createDoesAppExist(okta.NewAutoLoginApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "sign_on_url", "https://example.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "sign_on_redirect_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "reveal_password", "true"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "EDIT_USERNAME_AND_PASSWORD"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "sign_on_url", "https://exampleupdate.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "sign_on_redirect_url", "https://exampleupdate.com"),
					resource.TestCheckResourceAttr(resourceName, "reveal_password", "false"),
					resource.TestCheckResourceAttr(resourceName, "shared_password", "sharedpassword"),
					resource.TestCheckResourceAttr(resourceName, "shared_username", "sharedusername"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "SHARED_USERNAME_AND_PASSWORD"),
				),
			},
		},
	})
}

func buildTestAutoLoginConfigPreconfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  preconfigured_app		= "aws_console"
  label         		= "%s"
}
`, autoLoginApp, name, name)
}

func buildTestAutoLoginConfigPreconfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  preconfigured_app		= "aws_console"
  label         		= "%s"
  status 	   	 		= "INACTIVE"
}
`, autoLoginApp, name, name)
}

func buildTestAutoLoginConfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label         	 	= "%s"
  sign_on_url			= "https://example.com/login.html"
  sign_on_redirect_url	= "https://example.com"
  reveal_password		= true
  credentials_scheme 	= "EDIT_USERNAME_AND_PASSWORD"
}
`, autoLoginApp, name, name)
}

func buildTestAutoLoginConfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       			= "%s"
  status 	  			= "INACTIVE"
  sign_on_url			= "https://exampleupdate.com/login.html"
  sign_on_redirect_url	= "https://exampleupdate.com"
  reveal_password		= false
  credentials_scheme 	= "SHARED_USERNAME_AND_PASSWORD"
  shared_username 		= "sharedusername"
  shared_password		= "sharedpassword"
}
`, autoLoginApp, name, name)
}
