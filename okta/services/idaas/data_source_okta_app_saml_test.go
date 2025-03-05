package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccDataSourceOktaAppSaml_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppSaml, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	appCreate := buildTestAppSaml(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: appCreate,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_saml.test", "key_id"),
					resource.TestCheckResourceAttr("data.okta_app_saml.test", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_app_saml.test", "status", idaas.StatusActive),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "status", idaas.StatusActive),
					resource.TestCheckResourceAttr("data.okta_app_saml.test", "saml_signed_request_enabled", "false"),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "saml_signed_request_enabled", "false"),
				),
			},
		},
	})
}

func buildTestAppSaml(d int) string {
	return fmt.Sprintf(`
resource "okta_app_saml" "test" {
  label                    = "testAcc_%d"
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

  attribute_statements {
    type         = "GROUP"
    name         = "Attr Two"
    filter_type  = "STARTS_WITH"
    filter_value = "test"
  }
}`, d)
}
