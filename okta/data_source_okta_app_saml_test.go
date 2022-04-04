package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAppSaml_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSaml)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	appCreate := buildTestAppSaml(ri)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: appCreate,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_saml.test", "key_id"),
					resource.TestCheckResourceAttr("data.okta_app_saml.test", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr("data.okta_app_saml.test", "status", statusActive),
					resource.TestCheckResourceAttr("data.okta_app_saml.test_label", "status", statusActive),
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
