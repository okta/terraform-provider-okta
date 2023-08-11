package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaIdpSaml_crud(t *testing.T) {
	mgr := newFixtureManager(idpSaml, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", idpSaml)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(idpSaml, createDoesIdpExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "audience"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "sso_destination", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "sso_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_scope", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_scope", "ANY"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "max_clock_skew", "60"),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "audience"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "https://idp.example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "sso_destination", "https://idp.example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "sso_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "https://idp.example.com/issuer"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_scope", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_scope", "RESPONSE"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
		},
	})
}

// TestAccOktaIdpSaml_minimal_example was used to prove that the PR
// https://github.com/okta/terraform-provider-okta/pull/1355 was correct. This
// test would fail if the org was missing the mappings api feature. And pass if
// the feature was enabled.
func TestAccOktaIdpSaml_minimal_example(t *testing.T) {
	mgr := newFixtureManager(idpSaml, t.Name())
	config := `
resource "okta_app_saml" "test" {
	label                    = "testAcc_replace_with_uuid"
	sso_url                  = "http://google.com"
	recipient                = "http://here.com"
	destination              = "http://its-about-the-journey.com"
	audience                 = "http://audience.com"
	subject_name_id_template = "$${user.userName}"
	subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
	response_signed          = true
	signature_algorithm      = "RSA_SHA256"
	digest_algorithm         = "SHA256"
	honor_force_authn        = false
	authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}
resource "okta_idp_saml_key" "test" {
	x5c = [okta_app_saml.test.certificate]
}
resource "okta_idp_saml" "test" {
  name                     = "testAcc_replace_with_uuid"
  acs_type                 = "INSTANCE"
  sso_binding              = "HTTP-POST"
  sso_url                  = "https://idp.example.com"
  sso_destination          = "https://idp.example.com"
  username_template        = "idpuser.email"
  kid                      = okta_idp_saml_key.test.id
  issuer                   = "https://idp.example.com"
  request_signature_scope  = "REQUEST"
  response_signature_scope = "ANY"
}
	`
	resourceName := fmt.Sprintf("%s.test", idpSaml)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(idpSaml, createDoesIdpExist),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "audience"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "sso_destination", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "sso_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_scope", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_scope", "ANY"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
		},
	})
}
