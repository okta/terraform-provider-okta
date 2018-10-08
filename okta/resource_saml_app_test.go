package okta

import (
	"fmt"
	"regexp"
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

// Ensure conditional require logic causes this plan to fail
func TestAccOktaSamlApplicationConditionalRequire(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSamlConfigMissingFields(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(samlApp, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("missing conditionally required fields, reason: Custom SAML applications must contain these fields, missing fields: sso_url, recipient, destination, audience, idp_issuer, subject_name_id_template, subject_name_id_format, signature_algorithm, digest_algorithm, honor_force_authn, authn_context_class_ref"),
			},
		},
	})
}

// Ensure conditional require logic causes this plan to fail
func TestAccOktaSamlApplicationInvalidUrl(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSamlConfigInvalidUrl(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(samlApp, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("config is invalid: okta_saml_app.testAcc-.*: failed to validate url, \"123\""),
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
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
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
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "idp_issuer", "idhere123"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_template", "${source.login}"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"),
					resource.TestCheckResourceAttr(resourceName, "response_signed", "true"),
					resource.TestCheckResourceAttr(resourceName, "assertion_signed", "true"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "RSA_SHA1"),
					resource.TestCheckResourceAttr(resourceName, "digest_algorithm", "SHA1"),
					resource.TestCheckResourceAttr(resourceName, "honor_force_authn", "true"),
					resource.TestCheckResourceAttr(resourceName, "authn_context_class_ref", "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
				),
			},
		},
	})
}

func buildTestSamlConfigPreconfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  preconfigured_app	    = "amazon_aws"
  label         		= "%s"
}
`, samlApp, name, name)
}

func buildTestSamlConfigPreconfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  preconfigured_app	    = "amazon_aws"
  label         		= "%s"
  status 	    	    = "INACTIVE"
}
`, samlApp, name, name)
}

func buildTestSamlConfigMissingFields(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label         		= "%s"
  status 	    	    = "INACTIVE"
}
`, samlApp, name, name)
}

func buildTestSamlConfigInvalidUrl(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label         		= "%s"
  status 	    	    = "INACTIVE"
  sso_url      			= "123"
}
`, samlApp, name, name)
}

func buildTestSamlConfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       				= "%s"
  sso_url      				= "http://google.com"
  recipient 				= "http://here.com"
  destination 				= "http://its-about-the-journey.com"
  audience 					= "http://audience.com"
  idp_issuer 				= "idhere123"
  subject_name_id_template  = "$${user.userName}"
  subject_name_id_format	= "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed 			= true
  signature_algorithm 		= "RSA_SHA256"
  digest_algorithm 			= "SHA256"
  honor_force_authn			= false
  authn_context_class_ref	= "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}
`, samlApp, name, name)
}

func buildTestSamlConfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       = "%s"
  sso_url      				= "http://google.com"
  recipient 				= "http://here.com"
  destination 				= "http://its-about-the-journey.com"
  audience 					= "http://audience.com"
  idp_issuer 				= "idhere123"
  status 	  			    = "INACTIVE"
  subject_name_id_template  = "$${user.userName}"
  subject_name_id_format	= "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  signature_algorithm 		= "RSA_SHA256"
  response_signed 			= true
  digest_algorithm 			= "SHA256"
  honor_force_authn			= false
  authn_context_class_ref	= "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}
`, samlApp, name, name)
}

func buildTestSamlConfigAllFields(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       				= "%s"
  sso_url      				= "http://google.com"
  recipient 				= "http://here.com"
  destination 				= "http://its-about-the-journey.com"
  audience 					= "http://audience.com"
  idp_issuer 				= "idhere123"
  subject_name_id_template 	= "$${source.login}"
  subject_name_id_format	= "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
  response_signed 			= true
  assertion_signed 			= true
  signature_algorithm 		= "RSA_SHA1"
  digest_algorithm 			= "SHA1"
  honor_force_authn			= true
  authn_context_class_ref 	= "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}
`, samlApp, name, name)
}
