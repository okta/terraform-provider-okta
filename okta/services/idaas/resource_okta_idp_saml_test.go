package idaas_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaIdpSaml_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSIdpSaml, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSIdpSaml)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSIdpSaml, createDoesIdpExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
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
					resource.TestCheckResourceAttr(resourceName, "name_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
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
					resource.TestCheckResourceAttr(resourceName, "name_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
			{
				// Before fixing
				// https://github.com/okta/terraform-provider-okta/issues/1558
				// Not all settable arguments that were from API values were
				// being set on the read like sso_url.
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import resource into state")
					}
					expectedAttrs := []string{
						"acs_binding",
						"acs_type",
						"audience",
						"deprovisioned_action",
						"issuer",
						// "issuer_mode", not set during test
						"kid",
						"max_clock_skew",
						"name",
						"name_format",
						"profile_master",
						"provisioning_action",
						"sso_binding",
						"sso_destination",
						"sso_url",
						"status",
						// "subject_filter", not set during test
						// "subject_match_attribute", not set durting test
						"subject_match_type",
						"suspended_action",
						"user_type_id",
						"username_template",
					}
					notFound := []string{}
					for _, attr := range expectedAttrs {
						if s[0].Attributes[attr] == "" {
							notFound = append(notFound, attr)
						}
					}
					if len(notFound) > 0 {
						return fmt.Errorf("expected attributes %s to be set during import read", strings.Join(notFound, ", "))
					}
					return nil
				},
			},
		},
	})
}

// TestAccResourceOktaIdpSaml_minimal_example was used to prove that the PR
// https://github.com/okta/terraform-provider-okta/pull/1355 was correct. This
// test would fail if the org was missing the mappings api feature. And pass if
// the feature was enabled.
func TestAccResourceOktaIdpSaml_minimal_example(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSIdpSaml, t.Name())
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
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSIdpSaml)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSIdpSaml, createDoesIdpExist),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
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
