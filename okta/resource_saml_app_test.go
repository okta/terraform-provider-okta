package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

// Test creation of a simple AWS app. The preconfigured apps are created by name.
func TestAccOktaSamlApplicationPreconfig(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSamlConfigPreconfig(ri)
	updatedConfig := buildTestSamlConfigPreconfigUpdated(ri)
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
func TestAccOktaSamlApplication(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSamlConfig(ri)
	updatedConfig := buildTestSamlConfigUpdated(ri)
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

func TestAccOktaSamlApplicationAllFields(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSamlConfig(ri)
	allFields := buildTestSamlConfigAllFields(ri)
	resourceName := buildResourceFQN(samlApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(samlApp, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: allFields,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func buildTestSamlConfigPreconfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name		    = "amazon_aws"
  label         = "%s"
}
`, samlApp, name, name)
}

func buildTestSamlConfigPreconfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name		    = "amazon_aws"
  label         = "%s"
  status 	    = "INACTIVE"
}
`, samlApp, name, name)
}

func buildTestSamlConfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       = "%s"
  sso_url      				= "http://google.com"
  sso_url_override 			= "http://override.com/test"
  recipient 				= "http://here.com"
  recipient_override 		= "http://no-here.com"
  destination 				= "http://its-about-the-journey.com"
  destination_override 		= "http://out-of-order.com"
  audience 					= "http://audience.com"
  audience_override	 		= "http://stuff.com"
  idp_issuer 				= "idhere123"
}
`, samlApp, name, name)
}

func buildTestSamlConfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       = "%s"
  sso_url      				= "http://google.com"
  sso_url_override 			= "http://override.com/test"
  recipient 				= "http://here.com"
  recipient_override 		= "http://no-here.com"
  destination 				= "http://its-about-the-journey.com"
  destination_override 		= "http://out-of-order.com"
  audience 					= "http://audience.com"
  audience_override	 		= "http://stuff.com"
  idp_issuer 				= "idhere123"
  status 	  = "INACTIVE"
}
`, samlApp, name, name)
}

func buildTestSamlConfigAllFields(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       				= "%s"
  sso_url      				= "http://google.com"
  sso_url_override 			= "http://override.com/test"
  recipient 				= "http://here.com"
  recipient_override 		= "http://no-here.com"
  destination 				= "http://its-about-the-journey.com"
  destination_override 		= "http://out-of-order.com"
  audience 					= "http://audience.com"
  audience_override	 		= "http://stuff.com"
  idp_issuer 				= "idhere123"
  subject_name_id_template 	= "${fn:substringBefore(source.login, \"@\")}"
  subject_name_id_format 	= "EmailAddress"
  response_signed 			= true
  assertion_signed 			= true
  signature_algorithm 		= "RSA-SHA1"
  digest_algorithm 			= "SHA1"
  honor_force_authn			= true
  authn_context_class_ref 	= "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}
`, samlApp, name, name)
}
