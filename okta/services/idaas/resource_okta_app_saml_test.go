package idaas_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/sdk"
)

// Ensure conditional require logic causes this plan to fail
func TestAccResourceOktaAppSaml_conditionalRequire(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	config := buildTestSamlConfigMissingFields(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("missing conditionally required fields, reason: 'Custom SAML applications must contain these fields"),
			},
		},
	})
}

// Ensure conditional require logic causes this plan to fail
func TestAccResourceOktaAppSaml_invalidURL(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	config := buildTestSamlConfigMissingFields(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Custom SAML applications must contain these fields"),
			},
		},
	})
}

func TestAccResourceOktaAppSaml_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	allFields := mgr.GetFixtures("updated.tf", t)
	importConfig := mgr.GetFixtures("import.tf", t)
	importSAML11Config := mgr.GetFixtures("import_saml_1_1.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: allFields,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
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
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.1.type", "GROUP"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.1.filter_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.1.filter_value", "test"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.0", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.1", "https://okta.com"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "metadata_url"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.name", "groups"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.type", "GROUP"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.filter_type", "REGEX"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.filter_value", ".*"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_issuer", "https://dunshire.okta.com"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_url", "https://dunshire.okta.com/logout"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_certificate", "MIIFnDCCA4QCCQDBSLbiON2T1zANBgkqhkiG9w0BAQsFADCBjzELMAkGA1UEBhMCVVMxDjAMBgNV\r\nBAgMBU1haW5lMRAwDgYDVQQHDAdDYXJpYm91MRcwFQYDVQQKDA5Tbm93bWFrZXJzIEluYzEUMBIG\r\nA1UECwwLRW5naW5lZXJpbmcxDTALBgNVBAMMBFNub3cxIDAeBgkqhkiG9w0BCQEWEWVtYWlsQGV4\r\nYW1wbGUuY29tMB4XDTIwMTIwMzIyNDY0M1oXDTMwMTIwMTIyNDY0M1owgY8xCzAJBgNVBAYTAlVT\r\nMQ4wDAYDVQQIDAVNYWluZTEQMA4GA1UEBwwHQ2FyaWJvdTEXMBUGA1UECgwOU25vd21ha2VycyBJ\r\nbmMxFDASBgNVBAsMC0VuZ2luZWVyaW5nMQ0wCwYDVQQDDARTbm93MSAwHgYJKoZIhvcNAQkBFhFl\r\nbWFpbEBleGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBANMmWDjXPdoa\r\nPyzIENqeY9njLan2FqCbQPSestWUUcb6NhDsJVGSQ7XR+ozQA5TaJzbP7cAJUj8vCcbqMZsgOQAu\r\nO/pzYyQEKptLmrGvPn7xkJ1A1xLkp2NY18cpDTeUPueJUoidZ9EJwEuyUZIktzxNNU1pA1lGijiu\r\n2XNxs9d9JR/hm3tCu9Im8qLVB4JtX80YUa6QtlRjWR/H8a373AYCOASdoB3c57fIPD8ATDNy2w/c\r\nfCVGiyKDMFB+GA/WTsZpOP3iohRp8ltAncSuzypcztb2iE+jijtTsiC9kUA2abAJqqpoCJubNShi\r\nVff4822czpziS44MV2guC9wANi8u3Uyl5MKsU95j01jzadKRP5S+2f0K+n8n4UoV9fnqZFyuGAKd\r\nCJi9K6NlSAP+TgPe/JP9FOSuxQOHWJfmdLHdJD+evoKi9E55sr5lRFK0xU1Fj5Ld7zjC0pXPhtJf\r\nsgjEZzD433AsHnRzvRT1KSNCPkLYomznZo5n9rWYgCQ8HcytlQDTesmKE+s05E/VSWNtH84XdDrt\r\nieXwfwhHfaABSu+WjZYxi9CXdFCSvXhsgufUcK4FbYAHl/ga/cJxZc52yFC7Pcq0u9O2BSCjYPdQ\r\nDAHs9dhT1RhwVLM8RmoAzgxyyzau0gxnAlgSBD9FMW6dXqIHIp8yAAg9cRXhYRTNAgMBAAEwDQYJ\r\nKoZIhvcNAQELBQADggIBADofEC1SvG8qa7pmKCjB/E9Sxhk3mvUO9Gq43xzwVb721Ng3VYf4vGU3\r\nwLUwJeLt0wggnj26NJweN5T3q9T8UMxZhHSWvttEU3+S1nArRB0beti716HSlOCDx4wTmBu/D1MG\r\nt/kZYFJw+zuzvAcbYct2pK69AQhD8xAIbQvqADJI7cCK3yRry+aWtppc58P81KYabUlCfFXfhJ9E\r\nP72ffN4jVHpX3lxxYh7FKAdiKbY2FYzjsc7RdgKI1R3iAAZUCGBTvezNzaetGzTUjjl/g1tcVYij\r\nltH9ZOQBPlUMI88lxUxqgRTerpPmAJH00CACx4JFiZrweLM1trZyy06wNDQgLrqHr3EOagBF/O2h\r\nhfTehNdVr6iq3YhKWBo4/+RL0RCzHMh4u86VbDDnDn4Y6HzLuyIAtBFoikoKM6UHTOa0Pqv2bBr5\r\nwbkRkVUxl9yJJw/HmTCdfnsM9dTOJUKzEglnGF2184Gg+qJDZB6fSf0EAO1F6sTqiSswl+uHQZiy\r\nDaZzyU7Gg5seKOZ20zTRaX3Ihj9Zij/ORnrARE7eM/usKMECp+7syUwAUKxDCZkGiUdskmOhhBGL\r\nJtbyK3F2UvoJoLsm3pIcvMak9KwMjSTGJB47ABUP1+w+zGcNk0D5Co3IJ6QekiLfWJyQ+kKsWLKt\r\nzOYQQatrnBagM7MI2/T4\r\n"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: importConfig,
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import resource into state")
					}
					if s[0].Attributes["preconfigured_app"] != "pagerduty" {
						return errors.New("failed to set required properties when import existing infrastructure")
					}
					return nil
				},
			},
			{
				Config: importSAML11Config,
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import resource into state")
					}
					if s[0].Attributes["preconfigured_app"] != "sharepoint_onpremise" {
						return errors.New("failed to set required properties when import existing infrastructure")
					}
					return nil
				},
			},
		},
	})
}

func TestAccResourceOktaAppSaml_preconfigured(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	preconfigured := mgr.GetFixtures("preconfigured.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               wsFedAutoSSOErrorCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: preconfigured,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "preconfigured_app", "office365"),
					resource.TestCheckResourceAttr(resourceName, "saml_version", "1.1"),
					testAppSamlJson(resourceName, `{
       "wsFedConfigureType": "AUTO",
       "windowsTransportEnabled": false,
       "domain": "okta.com",
       "msftTenant": "okta",
       "domains": [],
       "requireAdminConsent": false
    }`, `{
      "calendar": false,
      "crm": false,
      "delve": false,
      "excel": false,
      "forms": false,
      "mail": false,
      "newsfeed": false,
      "onedrive": false,
      "people": false,
      "planner": false,
      "powerbi": false,
      "powerpoint": false,
      "sites": false,
      "sway": false,
      "tasks": false,
      "teams": false,
      "word": false,
      "yammer": false,
      "login": true
	}`),
				),
			},
		},
	})
}

func testAppSamlJson(name, expectedSettingsJSON, expectedLinksJSON string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		actualSettingsJSON := rs.Primary.Attributes["app_settings_json"]
		actualLinksJSON := rs.Primary.Attributes["app_links_json"]
		eq := areJSONStringsEqual(expectedSettingsJSON, actualSettingsJSON)
		if !eq {
			return fmt.Errorf("attribute 'app_settings_json' expected %q, got %q", expectedSettingsJSON, actualSettingsJSON)
		}
		eq = areJSONStringsEqual(expectedLinksJSON, actualLinksJSON)
		if !eq {
			return fmt.Errorf("attribute 'app_links_json' expected %q, got %q", expectedSettingsJSON, actualSettingsJSON)
		}
		return nil
	}
}

func areJSONStringsEqual(a, b string) bool {
	var aM, bM map[string]interface{}
	_ = json.Unmarshal([]byte(a), &aM)
	_ = json.Unmarshal([]byte(b), &bM)
	return reflect.DeepEqual(aM, bM)
}

func TestAccResourceOktaAppSaml_inlineHook(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	config := mgr.GetFixtures("basic_inline_hook.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet(resourceName, "inline_hook_id"),
				),
			},
		},
	})
}

// Tests creation of service app and updates it to turn on federated broker
func TestAccResourceOktaAppSaml_federationBroker(t *testing.T) {
	// TODO: This is an "Early Access Feature" and needs to be enabled by Okta
	//       Skipping for now assuming that the okta account doesn't have this feature enabled.
	//       If this feature is enabled or Okta releases this to all this test should be enabled.
	//       SEE https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm
	t.Skip("This is an 'Early Access Feature' and needs to be enabled by Okta, skipping this test as it fails when this feature is not available")

	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	config := mgr.GetFixtures("federation_broker_off.tf", t)
	updatedConfig := mgr.GetFixtures("federation_broker_on.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesAppExist(sdk.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "implicit_assignment", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "implicit_assignment", "true"),
				),
			},
		},
	})
}

func buildTestSamlConfigMissingFields(rInt int) string {
	name := acctest.BuildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label         		= "%s"
  status 	    	    = "INACTIVE"
}
`, resources.OktaIDaaSAppSaml, name, name)
}

func TestAccResourceOktaAppSaml_timeouts(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)
	config := `
resource "okta_app_saml" "test" {
  label                     = "testAcc_replace_with_uuid"
  sso_url                   = "http://google.com"
  recipient                 = "http://here.com"
  destination               = "http://its-about-the-journey.com"
  audience                  = "http://audience.com"
  subject_name_id_template  = "$${user.userName}"
  subject_name_id_format    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed           = true
  signature_algorithm       = "RSA_SHA256"
  digest_algorithm          = "SHA256"
  honor_force_authn         = false
  authn_context_class_ref   = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
  single_logout_issuer      = "https://dunshire.okta.com"
  single_logout_url         = "https://dunshire.okta.com/logout"
  single_logout_certificate = "MIIFnDCCA4QCCQDBSLbiON2T1zANBgkqhkiG9w0BAQsFADCBjzELMAkGA1UEBhMCVVMxDjAMBgNV\r\nBAgMBU1haW5lMRAwDgYDVQQHDAdDYXJpYm91MRcwFQYDVQQKDA5Tbm93bWFrZXJzIEluYzEUMBIG\r\nA1UECwwLRW5naW5lZXJpbmcxDTALBgNVBAMMBFNub3cxIDAeBgkqhkiG9w0BCQEWEWVtYWlsQGV4\r\nYW1wbGUuY29tMB4XDTIwMTIwMzIyNDY0M1oXDTMwMTIwMTIyNDY0M1owgY8xCzAJBgNVBAYTAlVT\r\nMQ4wDAYDVQQIDAVNYWluZTEQMA4GA1UEBwwHQ2FyaWJvdTEXMBUGA1UECgwOU25vd21ha2VycyBJ\r\nbmMxFDASBgNVBAsMC0VuZ2luZWVyaW5nMQ0wCwYDVQQDDARTbm93MSAwHgYJKoZIhvcNAQkBFhFl\r\nbWFpbEBleGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBANMmWDjXPdoa\r\nPyzIENqeY9njLan2FqCbQPSestWUUcb6NhDsJVGSQ7XR+ozQA5TaJzbP7cAJUj8vCcbqMZsgOQAu\r\nO/pzYyQEKptLmrGvPn7xkJ1A1xLkp2NY18cpDTeUPueJUoidZ9EJwEuyUZIktzxNNU1pA1lGijiu\r\n2XNxs9d9JR/hm3tCu9Im8qLVB4JtX80YUa6QtlRjWR/H8a373AYCOASdoB3c57fIPD8ATDNy2w/c\r\nfCVGiyKDMFB+GA/WTsZpOP3iohRp8ltAncSuzypcztb2iE+jijtTsiC9kUA2abAJqqpoCJubNShi\r\nVff4822czpziS44MV2guC9wANi8u3Uyl5MKsU95j01jzadKRP5S+2f0K+n8n4UoV9fnqZFyuGAKd\r\nCJi9K6NlSAP+TgPe/JP9FOSuxQOHWJfmdLHdJD+evoKi9E55sr5lRFK0xU1Fj5Ld7zjC0pXPhtJf\r\nsgjEZzD433AsHnRzvRT1KSNCPkLYomznZo5n9rWYgCQ8HcytlQDTesmKE+s05E/VSWNtH84XdDrt\r\nieXwfwhHfaABSu+WjZYxi9CXdFCSvXhsgufUcK4FbYAHl/ga/cJxZc52yFC7Pcq0u9O2BSCjYPdQ\r\nDAHs9dhT1RhwVLM8RmoAzgxyyzau0gxnAlgSBD9FMW6dXqIHIp8yAAg9cRXhYRTNAgMBAAEwDQYJ\r\nKoZIhvcNAQELBQADggIBADofEC1SvG8qa7pmKCjB/E9Sxhk3mvUO9Gq43xzwVb721Ng3VYf4vGU3\r\nwLUwJeLt0wggnj26NJweN5T3q9T8UMxZhHSWvttEU3+S1nArRB0beti716HSlOCDx4wTmBu/D1MG\r\nt/kZYFJw+zuzvAcbYct2pK69AQhD8xAIbQvqADJI7cCK3yRry+aWtppc58P81KYabUlCfFXfhJ9E\r\nP72ffN4jVHpX3lxxYh7FKAdiKbY2FYzjsc7RdgKI1R3iAAZUCGBTvezNzaetGzTUjjl/g1tcVYij\r\nltH9ZOQBPlUMI88lxUxqgRTerpPmAJH00CACx4JFiZrweLM1trZyy06wNDQgLrqHr3EOagBF/O2h\r\nhfTehNdVr6iq3YhKWBo4/+RL0RCzHMh4u86VbDDnDn4Y6HzLuyIAtBFoikoKM6UHTOa0Pqv2bBr5\r\nwbkRkVUxl9yJJw/HmTCdfnsM9dTOJUKzEglnGF2184Gg+qJDZB6fSf0EAO1F6sTqiSswl+uHQZiy\r\nDaZzyU7Gg5seKOZ20zTRaX3Ihj9Zij/ORnrARE7eM/usKMECp+7syUwAUKxDCZkGiUdskmOhhBGL\r\nJtbyK3F2UvoJoLsm3pIcvMak9KwMjSTGJB47ABUP1+w+zGcNk0D5Co3IJ6QekiLfWJyQ+kKsWLKt\r\nzOYQQatrnBagM7MI2/T4\r\n"
  attribute_statements {
    type         = "GROUP"
    name         = "groups"
    filter_type  = "REGEX"
    filter_value = ".*"
  }
  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}
`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}

// Test to ensure that certificate logic returns no-op / no-change upon apply and future plans
func TestAccResourceOktaAppSaml_certdiff(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	config := mgr.GetFixtures("basic_cert_plain.tf", t)
	config2 := mgr.GetFixtures("basic_cert_file.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.name", "groups"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.type", "GROUP"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.filter_type", "REGEX"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.filter_value", ".*"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_issuer", "https://dunshire.okta.com"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_url", "https://dunshire.okta.com/logout"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_certificate", "MIIFnDCCA4QCCQDBSLbiON2T1zANBgkqhkiG9w0BAQsFADCBjzELMAkGA1UEBhMCVVMxDjAMBgNV\r\nBAgMBU1haW5lMRAwDgYDVQQHDAdDYXJpYm91MRcwFQYDVQQKDA5Tbm93bWfrZXJzIEluYzEUMBIG\r\nA1UECwwLRW5naW5lZXJpbmcxDTALBgNVBAMMBFNub3cxIDAeBgkqhkiG9w0BCQEWEWVtYWlsQGV4\r\nYW1wbGUuY29tMB4XDTIwMTIwMzIyNDY0M1oXDTMwMTIwMTIyNDY0M1owgY8xCzAJBgNVBAYTAlVT\r\nMQ4wDAYDVQQIDAVNYWluZTEQMA4GA1UEBwwHQ2FyaWJvdTEXMBUGA1UECgwOU25vd21ha2VycyBJ\r\nbmMxFDASBgNVBAsMC0VuZ2luZWVyaW5nMQ0wCwYDVQQDDARTbm93MSAwHgYJKoZIhvcNAQkBFhFl\r\nbWFpbEBleGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBANMmWDjXPdoa\r\nPyzIENqeY9njLan2FqCbQPSestWUUcb6NhDsJVGSQ7XR+ozQA5TaJzbP7cAJUj8vCcbqMZsgOQAu\r\nO/pzYyQEKptLmrGvPn7xkJ1A1xLkp2NY18cpDTeUPueJUoidZ9EJwEuyUZIktzxNNU1pA1lGijiu\r\n2XNxs9d9JR/hm3tCu9Im8qLVB4JtX80YUa6QtlRjWR/H8a373AYCOASdoB3c57fIPD8ATDNy2w/c\r\nfCVGiyKDMFB+GA/WTsZpOP3iohRp8ltAncSuzypcztb2iE+jijtTsiC9kUA2abAJqqpoCJubNShi\r\nVff4822czpziS44MV2guC9wANi8u3Uyl5MKsU95j01jzadKRP5S+2f0K+n8n4UoV9fnqZFyuGAKd\r\nCJi9K6NlSAP+TgPe/JP9FOSuxQOHWJfmdLHdJD+evoKi9E55sr5lRFK0xU1Fj5Ld7zjC0pXPhtJf\r\nsgjEZzD433AsHnRzvRT1KSNCPkLYomznZo5n9rWYgCQ8HcytlQDTesmKE+s05E/VSWNtH84XdDrt\r\nieXwfwhHfaABSu+WjZYxi9CXdFCSvXhsgufUcK4FbYAHl/ga/cJxZc52yFC7Pcq0u9O2BSCjYPdQ\r\nDAHs9dhT1RhwVLM8RmoAzgxyyzau0gxnAlgSBD9FMW6dXqIHIp8yAAg9cRXhYRTNAgMBAAEwDQYJ\r\nKoZIhvcNAQELBQADggIBADofEC1SvG8qa7pmKCjB/E9Sxhk3mvUO9Gq43xzwVb721Ng3VYf4vGU3\r\nwLUwJeLt0wggnj26NJweN5T3q9T8UMxZhHSWvttEU3+S1nArRB0beti716HSlOCDx4wTmBu/D1MG\r\nt/kZYFJw+zuzvAcbYct2pK69AQhD8xAIbQvqADJI7cCK3yRry+aWtppc58P81KYabUlCfFXfhJ9E\r\nP72ffN4jVHpX3lxxYh7FKAdiKbY2FYzjsc7RdgKI1R3iAAZUCGBTvezNzaetGzTUjjl/g1tcVYij\r\nltH9ZOQBPlUMI88lxUxqgRTerpPmAJH00CACx4JFiZrweLM1trZyy06wNDQgLrqHr3EOagBF/O2h\r\nhfTehNdVr6iq3YhKWBo4/+RL0RCzHMh4u86VbDDnDn4Y6HzLuyIAtBFoikoKM6UHTOa0Pqv2bBr5\r\nwbkRkVUxl9yJJw/HmTCdfnsM9dTOJUKzEglnGF2184Gg+qJDZB6fSf0EAO1F6sTqiSswl+uHQZiy\r\nDaZzyU7Gg5seKOZ20zTRaX3Ihj9Zij/ORnrARE7eM/usKMECp+7syUwAUKxDCZkGiUdskmOhhBGL\r\nJtbyK3F2UvoJoLsm3pIcvMak9KwMjSTGJB47ABUP1+w+zGcNk0D5Co3IJ6QekiLfWJyQ+kKsWLKt\r\nzOYQQatrnBagM7MI2/T4\r\n"),
					resource.TestCheckResourceAttr(resourceName, "saml_signed_request_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.name", "groups"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.type", "GROUP"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.filter_type", "REGEX"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.filter_value", ".*"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_issuer", "https://dunshire.okta.com"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_url", "https://dunshire.okta.com/logout"),
					resource.TestCheckResourceAttr(resourceName, "single_logout_certificate", "MIIFnDCCA4QCCQDBSLbiON2T1zANBgkqhkiG9w0BAQsFADCBjzELMAkGA1UEBhMCVVMxDjAMBgNV\r\nBAgMBU1haW5lMRAwDgYDVQQHDAdDYXJpYm91MRcwFQYDVQQKDA5Tbm93bWfrZXJzIEluYzEUMBIG\r\nA1UECwwLRW5naW5lZXJpbmcxDTALBgNVBAMMBFNub3cxIDAeBgkqhkiG9w0BCQEWEWVtYWlsQGV4\r\nYW1wbGUuY29tMB4XDTIwMTIwMzIyNDY0M1oXDTMwMTIwMTIyNDY0M1owgY8xCzAJBgNVBAYTAlVT\r\nMQ4wDAYDVQQIDAVNYWluZTEQMA4GA1UEBwwHQ2FyaWJvdTEXMBUGA1UECgwOU25vd21ha2VycyBJ\r\nbmMxFDASBgNVBAsMC0VuZ2luZWVyaW5nMQ0wCwYDVQQDDARTbm93MSAwHgYJKoZIhvcNAQkBFhFl\r\nbWFpbEBleGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBANMmWDjXPdoa\r\nPyzIENqeY9njLan2FqCbQPSestWUUcb6NhDsJVGSQ7XR+ozQA5TaJzbP7cAJUj8vCcbqMZsgOQAu\r\nO/pzYyQEKptLmrGvPn7xkJ1A1xLkp2NY18cpDTeUPueJUoidZ9EJwEuyUZIktzxNNU1pA1lGijiu\r\n2XNxs9d9JR/hm3tCu9Im8qLVB4JtX80YUa6QtlRjWR/H8a373AYCOASdoB3c57fIPD8ATDNy2w/c\r\nfCVGiyKDMFB+GA/WTsZpOP3iohRp8ltAncSuzypcztb2iE+jijtTsiC9kUA2abAJqqpoCJubNShi\r\nVff4822czpziS44MV2guC9wANi8u3Uyl5MKsU95j01jzadKRP5S+2f0K+n8n4UoV9fnqZFyuGAKd\r\nCJi9K6NlSAP+TgPe/JP9FOSuxQOHWJfmdLHdJD+evoKi9E55sr5lRFK0xU1Fj5Ld7zjC0pXPhtJf\r\nsgjEZzD433AsHnRzvRT1KSNCPkLYomznZo5n9rWYgCQ8HcytlQDTesmKE+s05E/VSWNtH84XdDrt\r\nieXwfwhHfaABSu+WjZYxi9CXdFCSvXhsgufUcK4FbYAHl/ga/cJxZc52yFC7Pcq0u9O2BSCjYPdQ\r\nDAHs9dhT1RhwVLM8RmoAzgxyyzau0gxnAlgSBD9FMW6dXqIHIp8yAAg9cRXhYRTNAgMBAAEwDQYJ\r\nKoZIhvcNAQELBQADggIBADofEC1SvG8qa7pmKCjB/E9Sxhk3mvUO9Gq43xzwVb721Ng3VYf4vGU3\r\nwLUwJeLt0wggnj26NJweN5T3q9T8UMxZhHSWvttEU3+S1nArRB0beti716HSlOCDx4wTmBu/D1MG\r\nt/kZYFJw+zuzvAcbYct2pK69AQhD8xAIbQvqADJI7cCK3yRry+aWtppc58P81KYabUlCfFXfhJ9E\r\nP72ffN4jVHpX3lxxYh7FKAdiKbY2FYzjsc7RdgKI1R3iAAZUCGBTvezNzaetGzTUjjl/g1tcVYij\r\nltH9ZOQBPlUMI88lxUxqgRTerpPmAJH00CACx4JFiZrweLM1trZyy06wNDQgLrqHr3EOagBF/O2h\r\nhfTehNdVr6iq3YhKWBo4/+RL0RCzHMh4u86VbDDnDn4Y6HzLuyIAtBFoikoKM6UHTOa0Pqv2bBr5\r\nwbkRkVUxl9yJJw/HmTCdfnsM9dTOJUKzEglnGF2184Gg+qJDZB6fSf0EAO1F6sTqiSswl+uHQZiy\r\nDaZzyU7Gg5seKOZ20zTRaX3Ihj9Zij/ORnrARE7eM/usKMECp+7syUwAUKxDCZkGiUdskmOhhBGL\r\nJtbyK3F2UvoJoLsm3pIcvMak9KwMjSTGJB47ABUP1+w+zGcNk0D5Co3IJ6QekiLfWJyQ+kKsWLKt\r\nzOYQQatrnBagM7MI2/T4\r\n"),
					resource.TestCheckResourceAttr(resourceName, "saml_signed_request_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config:   config2,
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceOktaAppSaml_Issue2021(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)
	config := `
resource "okta_app_saml" "test" {
  accessibility_self_service     = "false"
  assertion_signed               = "true"
  audience                       = "https://example.com/audience"
  authn_context_class_ref        = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
  auto_submit_toolbar            = "false"
  default_relay_state            = "/"
  destination                    = "https://example.com/audience"
  digest_algorithm               = "SHA256"
  hide_ios                       = "true"
  hide_web                       = "true"
  honor_force_authn              = "true"
  idp_issuer                     = "http://www.okta.com/$${org.externalKey}"
  implicit_assignment            = "false"
  label                          = "SAML APP"
  recipient                      = "https://example.com/audience"
  response_signed                = "true"
  saml_signed_request_enabled    = "true"
  saml_version                   = "2.0"
  signature_algorithm            = "RSA_SHA256"
  sso_url                        = "https://example.com/sso"
  status                         = "ACTIVE"
  subject_name_id_format         = "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent"
  subject_name_id_template       = "$${user.userName}"
  user_name_template             = "user.login"
  user_name_template_push_status = "PUSH"
  user_name_template_type        = "CUSTOM"
  single_logout_certificate      = "MIID2zCCAsOgAwIBAgIUHBaBGrGVVkp2kC+yPrhXc5N2+4swDQYJKoZIhvcNAQEL\r\nBQAwfTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM\r\nGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEUMBIGA1UEAwwLZXhhbXBsZS5jb20x\r\nIDAeBgkqhkiG9w0BCQEWEWhlbGxvQGV4YW1wbGUuY29tMB4XDTI0MDYxODAzNTU0\r\nMVoXDTI1MDYxODAzNTU0MVowfTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUt\r\nU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEUMBIGA1UE\r\nAwwLZXhhbXBsZS5jb20xIDAeBgkqhkiG9w0BCQEWEWhlbGxvQGV4YW1wbGUuY29t\r\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4G7iJpERa657fZdWVKpM\r\nxY+8KBtTe/bPx7v+7ccOA9JhsGoiJIilaqTEGi+VmLS0yBJJ75e0eRuCufXxdUU9\r\ncPtze6vVppIXjNDYKkCb4FpMJCXDR94ojYD28Q4j7R+A5MgoVaL4m6bQMxN4Gtu4\r\nww9tVoXXMtKlYm57Z+44KZ9zX9ZT7h5tpPk4bws2ooi3mv8tpPhh63s+eSdShL/0\r\nPLcaTRmeL3tCZ2R07Ea7ZHZix+DSAFGZ3MfhE0/q8PoEj8WSuvJtL7XhRq1xUsFL\r\nEQGGZNy4DJecu6mjhieKpsaQGSpMrMcmekvLaEtL6bOepDqVBsyzyvCzM+46LXGd\r\nhQIDAQABo1MwUTAdBgNVHQ4EFgQUlyqz0r9lJLuXVGY6XocwikJMzfIwHwYDVR0j\r\nBBgwFoAUlyqz0r9lJLuXVGY6XocwikJMzfIwDwYDVR0TAQH/BAUwAwEB/zANBgkq\r\nhkiG9w0BAQsFAAOCAQEAo9aqKVV+zIpaosBxCN5GQIhY6soa8FgEhcZrZvd2iL67\r\n9aLYDY46RnJgpa4RS+M0gTlp9u+3dH6uvuo8CmR243IOGH9LOWd624UN+tka+3PM\r\n50A7Uxo3KFfmOZi+ym5xn+UADJx8uUrH1owlMhFZMPWLr/JuoBAxVNI8KRXFhW4U\r\npcHmKvqU7GZo7m2QwE0JIJ5p00ED66jNky/IAqoexikbhZ8IgzTbtlWFzbqVKNq1\r\nzvcCEc4LXKytMQCCWv71HBNMfBvR4tEbcKmxe356IHcs+dmEFtg3dfEBfH5U5VoS\r\n1RqP+9+AB4coGpnm7F660PSwfyQwBZo5/a0HLqbZFA=="
}

`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttrSet(resourceName, "single_logout_certificate"),
					resource.TestCheckNoResourceAttr(resourceName, "single_logout_url"),
					resource.TestCheckNoResourceAttr(resourceName, "single_logout_issuer"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSaml_Issue2171AcsEndpoints(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)
	config := `
resource "okta_app_saml" "test" {
  acs_endpoints = [%s]

  label                    = "testAcc_replace_with_uuid"
  sso_url                  = "http://google.com"
  recipient                = "http://here.com"
  destination              = "http://its-about-the-journey.com"
  audience                 = "http://audience.com"
  subject_name_id_template = "$${source.login}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
  response_signed          = true
  assertion_signed         = true
  signature_algorithm      = "RSA_SHA1"
  digest_algorithm         = "SHA1"
  honor_force_authn        = true
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}
`
	acsEndpoints1 := "\"https://example.com\",\"https://okta.com\""
	acsEndpoints2 := "\"https://okta.com\",\"https://example.com\""
	acsEndpoints3 := "\"https://okta.com\",\"https://middle.example.com\",\"https://example.com\""

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(config, acsEndpoints1)),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.0", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.1", "https://okta.com"),
				),
			},
			{
				// demonstrate flipping order is respected
				Config: mgr.ConfigReplace(fmt.Sprintf(config, acsEndpoints2)),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.0", "https://okta.com"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.1", "https://example.com"),
				),
			},
			{
				// demonstrate inserting and order is respected
				Config: mgr.ConfigReplace(fmt.Sprintf(config, acsEndpoints3)),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.0", "https://okta.com"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.1", "https://middle.example.com"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints.2", "https://example.com"),
				),
			},
		},
	})
}

func wsFedAutoSSOErrorCheck(t *testing.T) resource.ErrorCheckFunc {
	return func(err error) error {
		skip, _ := regexp.MatchString("In order to continue using WS-FED Auto for SSO", err.Error())
		if skip {
			t.Skipf("must grant admin consent for WS-FED o365 app: %+v", err)
			return nil
		}
		return err
	}
}

func TestAccResourceOktaAppSaml_Issue2171AcsEndpointsWithIndex(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	baseConfig := mgr.GetFixtures("resource_acs_endpoints_indices.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints_indices.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints_indices.0.url", "https://example2.com"),
					resource.TestCheckResourceAttr(resourceName, "acs_endpoints_indices.0.index", "102"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSaml_skipAuthenticationPolicy(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)
	config := `
resource "okta_app_saml" "test" {
	label                    = "testAcc_replace_with_uuid"
	sso_url                  = "http://google.com"
	recipient                = "http://here.com"
	destination              = "http://its-about-the-journey.com"
	audience                 = "http://audience.com"
	subject_name_id_template = "$${source.login}"
	subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	response_signed          = true
	assertion_signed         = true
	signature_algorithm      = "RSA_SHA1"
	digest_algorithm         = "SHA1"
	honor_force_authn        = true
	authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
	skip_authentication_policy = true
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "skip_authentication_policy", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "metadata_url"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSaml_skipAuthenticationPolicyUpdate(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)
	config := `
resource "okta_app_saml" "test" {
	label                    = "testAcc_replace_with_uuid"
	sso_url                  = "http://google.com"
	recipient                = "http://here.com"
	destination              = "http://its-about-the-journey.com"
	audience                 = "http://audience.com"
	subject_name_id_template = "$${source.login}"
	subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	response_signed          = true
	assertion_signed         = true
	signature_algorithm      = "RSA_SHA1"
	digest_algorithm         = "SHA1"
	honor_force_authn        = true
	authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
	skip_authentication_policy = false
}`
	updatedConfig := `
resource "okta_app_saml" "test" {
	label                    = "testAcc_replace_with_uuid"
	sso_url                  = "http://google.com"
	recipient                = "http://here.com"
	destination              = "http://its-about-the-journey.com"
	audience                 = "http://audience.com"
	subject_name_id_template = "$${source.login}"
	subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	response_signed          = true
	assertion_signed         = true
	signature_algorithm      = "RSA_SHA1"
	digest_algorithm         = "SHA1"
	honor_force_authn        = true
	authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
	skip_authentication_policy = true
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "skip_authentication_policy", "false"),
				),
			},
			{
				Config: mgr.ConfigReplace(updatedConfig),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "skip_authentication_policy", "true"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSaml_skipAuthenticationPolicyOffice365(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSaml, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSaml)
	config := `
resource "okta_app_saml" "test" {
	label                    = "office365"
	sso_url                  = "http://google.com"
	recipient                = "http://here.com"
	destination              = "http://its-about-the-journey.com"
	audience                 = "http://audience.com"
	subject_name_id_template = "$${source.login}"
	subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	response_signed          = true
	assertion_signed         = true
	signature_algorithm      = "RSA_SHA1"
	digest_algorithm         = "SHA1"
	honor_force_authn        = true
	authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
	skip_authentication_policy = true
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", "office365"),
					resource.TestCheckResourceAttr(resourceName, "skip_authentication_policy", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "metadata_url"),
				),
			},
		},
	})
}
