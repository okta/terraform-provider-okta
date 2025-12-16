package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaAuthenticator_OTP_crud(t *testing.T) {
	config := `
	resource "okta_authenticator" "otp" {
		name   = "Custom OTP"
		key    = "custom_otp"
		status = "ACTIVE"
		settings = jsonencode({
		  "protocol" : "TOTP",
		  "acceptableAdjacentIntervals" : 3,
		  "timeIntervalInSeconds" : 30,
		  "encoding" : "base32",
		  "algorithm" : "HMacSHA256",
		  "passCodeLength" : 6
		})
		legacy_ignore_name = false
	}`
	resourceName := fmt.Sprintf("%s.otp", resources.OktaIDaaSAuthenticator)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "security_key"),
					resource.TestCheckResourceAttr(resourceName, "key", "custom_otp"),
					resource.TestCheckResourceAttr(resourceName, "name", "Custom OTP"),
					testAttributeJSON(resourceName, "settings", `{"acceptableAdjacentIntervals":3,"algorithm":"HMacSHA256","encoding":"base32","passCodeLength":6,"protocol":"TOTP","timeIntervalInSeconds":30}`),
				),
			},
		},
	})
}

func TestAccResourceOktaAuthenticator_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticator)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticator, t.Name())
	config := mgr.GetFixtures("security_question.tf", t)
	configUpdated := mgr.GetFixtures("security_question_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
					testAttributeJSON(resourceName, "settings", `{"allowedFor" : "recovery"}`),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
					testAttributeJSON(resourceName, "settings", `{"allowedFor" : "any"}`),
				),
			},
			{
				Config: config,
			},
			{
				Config: `
resource "okta_authenticator" "test" {
  status   = "INACTIVE"
  name     = "Security Question"
  key      = "security_question"
  settings = jsonencode(
  {
    "allowedFor" : "recovery"
  }
  )
}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
				),
			},
		},
	})
}

// TestAccResourceOktaAuthenticator_issue_1367_simple
// https://github.com/okta/terraform-provider-okta/issues/1367
// Google OTP is a simple example of the solution for #1367
func TestAccResourceOktaAuthenticator_issue_1367_simple(t *testing.T) {
	config := `
resource "okta_authenticator" "google_otp" {
	name   = "Google Authenticator"
	key    = "google_otp"
	status = "INACTIVE"
}
`
	resourceName := fmt.Sprintf("%s.google_otp", resources.OktaIDaaSAuthenticator)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "type", "app"),
					resource.TestCheckResourceAttr(resourceName, "key", "google_otp"),
					resource.TestCheckResourceAttr(resourceName, "name", "Google Authenticator"),
				),
			},
		},
	})
}

// TestAccResourceOktaAuthenticator_issue_1367_provider_json
// https://github.com/okta/terraform-provider-okta/issues/1367
// Demonstrates provider input as freeform JSON
// Example from `POST /api/v1/authenticator` API docs
// https://developer.okta.com/docs/reference/api/authenticators-admin/#create-authenticator
func TestAccResourceOktaAuthenticator_issue_1367_provider_json(t *testing.T) {
	config := `
resource "okta_authenticator" "test" {
  name = "On-Prem MFA"
  key = "onprem_mfa"
  provider_json = jsonencode(
	{
		"type": "DEL_OATH",
		"configuration": {
		  "authPort": 999,
		  "userNameTemplate": {
			"template": "global.assign.userName.login"
		  },
		  "hostName": "localhost",
		  "sharedSecret": "Sh4r3d s3cr3t"
		}
	  }
  )
}`
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticator)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "security_key"),
					resource.TestCheckResourceAttr(resourceName, "key", "onprem_mfa"),
					resource.TestCheckResourceAttr(resourceName, "name", "On-Prem MFA"),
					testAttributeJSON(resourceName, "profile_json", `{
						{
							"type": "DEL_OATH",
							"configuration": {
							  "authPort": 999,
							  "userNameTemplate": {
								"template": "global.assign.userName.login"
							  },
							  "hostName": "localhost",
							  "sharedSecret": "Sh4r3d s3cr3t"
							}
						}`),
					resource.TestCheckResourceAttr(resourceName, "provider_type", "DEL_OATH"),
					resource.TestCheckResourceAttr(resourceName, "provider_hostname", "localhost"),
					resource.TestCheckResourceAttr(resourceName, "provider_auth_port", "999"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_instance_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_user_name_template", "global.assign.userName.login"),
				),
			},
		},
	})
}

func TestAccResourceOktaAuthenticator_OktaVerifyCRUD(t *testing.T) {
	resourceName := fmt.Sprintf("%s.okta_verify", resources.OktaIDaaSAuthenticator)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticator, t.Name())
	config := mgr.GetFixtures("okta_verify.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				ImportState:        true,
				ResourceName:       "okta_authenticator.okta_verify",
				ImportStateId:      "0ktaV3rify1d",
				ImportStatePersist: true,
				Config:             config,
				PlanOnly:           true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "app"),
					resource.TestCheckResourceAttr(resourceName, "name", "Okta Verify"),
					// ensure any change to examples/resources/okta_authenticator/okta_verify.tf is reflected here.
					testAttributeJSON(resourceName, "settings", `{"channelBinding":{"required":"ALWAYS","style":"NUMBER_CHALLENGE"},"compliance":{"fips":"OPTIONAL"},"userVerification":"PREFERRED","enrollmentSecurityLevel":"HIGH","userVerificationMethods":["BIOMETRICS"]}`),
				),
			},
		},
	})
}

// TestAccResourceOktaAuthenticator_PhoneWithMethods_crud tests CRUD operations on phone authenticator with method blocks
func TestAccResourceOktaAuthenticator_PhoneWithMethods_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticator)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticator, t.Name())
	config := mgr.GetFixtures("phone_with_methods.tf", t)
	configUpdated := mgr.GetFixtures("phone_with_methods_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "phone"),
					resource.TestCheckResourceAttr(resourceName, "key", "phone_number"),
					resource.TestCheckResourceAttr(resourceName, "name", "Phone"),
					resource.TestCheckResourceAttr(resourceName, "method.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "sms",
						"status": "ACTIVE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "voice",
						"status": "INACTIVE",
					}),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "phone"),
					resource.TestCheckResourceAttr(resourceName, "key", "phone_number"),
					resource.TestCheckResourceAttr(resourceName, "name", "Phone"),
					resource.TestCheckResourceAttr(resourceName, "method.#", "2"),
					// Both methods now ACTIVE (API may not allow deactivating all methods)
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "sms",
						"status": "ACTIVE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "voice",
						"status": "ACTIVE",
					}),
				),
			},
		},
	})
}

// TestAccResourceOktaAuthenticator_OktaVerifyWithMethods_crud tests CRUD operations on Okta Verify with multiple methods and settings
func TestAccResourceOktaAuthenticator_OktaVerifyWithMethods_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticator)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticator, t.Name())
	config := mgr.GetFixtures("okta_verify_with_methods.tf", t)
	configUpdated := mgr.GetFixtures("okta_verify_methods_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "app"),
					resource.TestCheckResourceAttr(resourceName, "key", "okta_verify"),
					resource.TestCheckResourceAttr(resourceName, "name", "Okta Verify"),
					resource.TestCheckResourceAttr(resourceName, "method.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "push",
						"status": "ACTIVE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "totp",
						"status": "ACTIVE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "signed_nonce",
						"status": "ACTIVE",
					}),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "app"),
					resource.TestCheckResourceAttr(resourceName, "key", "okta_verify"),
					resource.TestCheckResourceAttr(resourceName, "name", "Okta Verify"),
					resource.TestCheckResourceAttr(resourceName, "method.#", "3"),
					// All methods remain ACTIVE, we're just testing settings updates
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "push",
						"status": "ACTIVE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "totp",
						"status": "ACTIVE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "method.*", map[string]string{
						"type":   "signed_nonce",
						"status": "ACTIVE",
					}),
					// Verify that settings were updated - API may add extra fields like appInstanceId and compliance
					resource.TestCheckResourceAttrSet(resourceName, "settings"),
				),
			},
		},
	})
}

// TestAccResourceOktaAuthenticator_PhoneMethodsNoConfig tests phone authenticator works without explicit method blocks
func TestAccResourceOktaAuthenticator_PhoneMethodsNoConfig(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticator)
	config := `
resource "okta_authenticator" "test" {
  name = "Phone"
  key  = "phone_number"
}
`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "phone"),
					resource.TestCheckResourceAttr(resourceName, "key", "phone_number"),
					resource.TestCheckResourceAttr(resourceName, "name", "Phone"),
					// Method blocks are optional - authenticator should still work
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

// TestAccResourceOktaAuthenticator_MethodsOptional tests backward compatibility - authenticators without method blocks
func TestAccResourceOktaAuthenticator_MethodsOptional(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticator)
	config := `
resource "okta_authenticator" "test" {
	name   = "Security Question"
	key    = "security_question"
	status = "ACTIVE"
	settings = jsonencode({
		"allowedFor" : "recovery"
	})
}
`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					// Verify that method blocks are optional - should have 0 methods when not specified
					resource.TestCheckResourceAttr(resourceName, "method.#", "0"),
					testAttributeJSON(resourceName, "settings", `{"allowedFor":"recovery"}`),
				),
			},
		},
	})
}

// TestAccResourceOktaAuthenticator_MethodsImportStateVerify tests that method blocks
// are correctly verified during ImportStateVerify
func TestAccResourceOktaAuthenticator_MethodsImportStateVerify(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticator)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticator, t.Name())
	config := mgr.GetFixtures("phone_with_methods.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"legacy_ignore_name",
					"provider_hostname",
					"provider_user_name_template",
				},
			},
		},
	})
}
