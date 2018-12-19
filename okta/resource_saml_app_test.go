package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

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
	mgr := newFixtureManager(samlApp)
	config := mgr.GetFixtures("custom_saml_app.tf", ri, t)
	updatedConfig := mgr.GetFixtures("custom_saml_app_updated.tf", ri, t)
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
	mgr := newFixtureManager(samlApp)
	config := mgr.GetFixtures("custom_saml_app.tf", ri, t)
	allFields := mgr.GetFixtures("custom_saml_app_all_fields.tf", ri, t)
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
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.name", "Attr One"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.namespace", "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.values.0", "val"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.1.name", "Attr Two"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.1.namespace", "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"),
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

// Add and remove groups/users
func TestAccOktaSamlApplicationUserGroups(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(samlApp)
	config := mgr.GetFixtures("saml_app_with_groups_and_users.tf", ri, t)
	updatedConfig := mgr.GetFixtures("saml_app_with_groups_and_users_updated.tf", ri, t)
	resourceName := buildResourceFQN(samlApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(samlApp, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0"),
					resource.TestCheckResourceAttr(resourceName, "key.years_valid", "3"),
					// resource.TestCheckResourceAttr(resourceName, "features.#", "1"),
					// resource.TestCheckResourceAttr(resourceName, "features.0", "PUSH_NEW_USERS"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckNoResourceAttr(resourceName, "users.0"),
					resource.TestCheckNoResourceAttr(resourceName, "groups.0"),
					resource.TestCheckNoResourceAttr(resourceName, "key.id"),
				),
			},
		},
	})
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
