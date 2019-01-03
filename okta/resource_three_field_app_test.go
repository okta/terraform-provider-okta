package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

func TestAccOktaThreeFieldApplication(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestThreeFieldConfig(ri)
	updatedConfig := buildTestThreeFieldConfigUpdated(ri)
	resourceName := buildResourceFQN(threeFieldApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(threeFieldApp, createDoesAppExist(okta.NewSwaThreeFieldApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaThreeFieldApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "button_selector", "btn"),
					resource.TestCheckResourceAttr(resourceName, "username_selector", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_selector", "pass"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_selector", "third"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_value", "third"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://example.com"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaThreeFieldApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "button_selector", "btn1"),
					resource.TestCheckResourceAttr(resourceName, "username_selector", "user1"),
					resource.TestCheckResourceAttr(resourceName, "password_selector", "pass1"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://example.com"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_selector", "mfa"),
					resource.TestCheckResourceAttr(resourceName, "extra_field_value", "mfa"),
				),
			},
		},
	})
}

func buildTestThreeFieldConfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label         	 	= "%s"
  button_selector		= "btn"
  username_selector		= "user"
  password_selector		= "pass"
  url			        = "http://example.com"
  extra_field_selector          = "third"
  extra_field_value		= "third"
}
`, threeFieldApp, name, name)
}

func buildTestThreeFieldConfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       			= "%s"
  status 	  		= "INACTIVE"
  button_selector		= "btn1"
  username_selector		= "user1"
  password_selector		= "pass1"
  url			        = "http://example.com"
  extra_field_selector 	        = "mfa"
  extra_field_value		= "mfa"
}
`, threeFieldApp, name, name)
}
