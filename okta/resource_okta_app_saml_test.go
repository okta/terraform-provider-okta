package okta

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/sdk"
)

// Ensure conditional require logic causes this plan to fail
func TestAccAppSaml_conditionalRequire(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSamlConfigMissingFields(ri)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("missing conditionally required fields, reason: 'Custom SAML applications must contain these fields"),
			},
		},
	})
}

// Ensure conditional require logic causes this plan to fail
func TestAccAppSaml_invalidURL(t *testing.T) {
	mgr := newFixtureManager(appSaml, t.Name())
	config := buildTestSamlConfigMissingFields(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Custom SAML applications must contain these fields"),
			},
		},
	})
}

func TestAccAppSaml_crud(t *testing.T) {
	mgr := newFixtureManager(appSaml, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	allFields := mgr.GetFixtures("updated.tf", t)
	importConfig := mgr.GetFixtures("import.tf", t)
	importSAML11Config := mgr.GetFixtures("import_saml_1_1.tf", t)
	resourceName := fmt.Sprintf("%s.test", appSaml)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: allFields,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
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
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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

func TestAccAppSaml_preconfigured(t *testing.T) {
	mgr := newFixtureManager(appSaml, t.Name())
	preconfigured := mgr.GetFixtures("preconfigured.tf", t)
	resourceName := fmt.Sprintf("%s.test", appSaml)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: preconfigured,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
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

func TestAccAppSaml_inlineHook(t *testing.T) {
	mgr := newFixtureManager(appSaml, t.Name())
	config := mgr.GetFixtures("basic_inline_hook.tf", t)
	resourceName := fmt.Sprintf("%s.test", appSaml)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttrSet(resourceName, "inline_hook_id"),
				),
			},
		},
	})
}

// Tests creation of service app and updates it to turn on federated broker
func TestAccAppSaml_federationBroker(t *testing.T) {
	// TODO: This is an "Early Access Feature" and needs to be enabled by Okta
	//       Skipping for now assuming that the okta account doesn't have this feature enabled.
	//       If this feature is enabled or Okta releases this to all this test should be enabled.
	//       SEE https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm
	t.Skip("This is an 'Early Access Feature' and needs to be enabled by Okta, skipping this test as it fails when this feature is not available")

	mgr := newFixtureManager(appSaml, t.Name())
	config := mgr.GetFixtures("federation_broker_off.tf", t)
	updatedConfig := mgr.GetFixtures("federation_broker_on.tf", t)
	resourceName := fmt.Sprintf("%s.test", appSaml)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appOAuth, createDoesAppExist(sdk.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "implicit_assignment", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "implicit_assignment", "true"),
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
`, appSaml, name, name)
}

func TestAccResourceOktaAppSaml_timeouts(t *testing.T) {
	mgr := newFixtureManager(appSaml, t.Name())
	resourceName := fmt.Sprintf("%s.test", appSaml)
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
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}

// Test to ensure that certificate logic returns no-op / no-change upon apply and future plans
func TestAccAppSaml_certdiff(t *testing.T) {
	mgr := newFixtureManager(appSaml, t.Name())
	config := mgr.GetFixtures("basic_cert_plain.tf", t)
	config2 := mgr.GetFixtures("basic_cert_file.tf", t)
	resourceName := fmt.Sprintf("%s.test", appSaml)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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
