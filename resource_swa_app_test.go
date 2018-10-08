package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

// Test creation of a simple AWS SWA app. The preconfigured apps are created by name.
func TestAccOktaSwaApplicationPreconfig(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSwaConfigPreconfig(ri)
	updatedConfig := buildTestSwaConfigPreconfigUpdated(ri)
	resourceName := buildResourceFQN(samlApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(samlApp, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
		},
	})
}

// Test creation of a custom SAML app.
func TestAccOktaSwaApplication(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSwaConfig(ri)
	updatedConfig := buildTestSwaConfigUpdated(ri)
	resourceName := buildResourceFQN(samlApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(samlApp, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
		},
	})
}

func buildTestSwaConfigPreconfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name		    = "amazon_aws"
  label         = "%s"
}
`, samlApp, name, name)
}

func buildTestSwaConfigPreconfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name		    = "amazon_aws"
  label         = "%s"
  status 	    = "INACTIVE"
}
`, samlApp, name, name)
}

func buildTestSwaConfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       = "%s"
  sso_url      = "http://google.com"
}
`, samlApp, name, name)
}

func buildTestSwaConfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       = "%s"
  sso_url      = "http://google.com"
  status 	  = "INACTIVE"
}
`, samlApp, name, name)
}
