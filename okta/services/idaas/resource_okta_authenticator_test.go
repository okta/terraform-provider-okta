package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

// TestAccResourceOktaAuthenticator_custom_app_crud
// Tests custom_app type authenticator with create, read, update operations
// custom_app authenticators do not support deletion (only deactivation)
func TestAccResourceOktaAuthenticator_custom_app_crud(t *testing.T) {
	resourceName1 := fmt.Sprintf("%s.test1", resources.OktaIDaaSAuthenticator)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticator, t.Name())
	config := mgr.GetFixtures("custom_app.tf", t)
	updatedConfig := mgr.GetFixtures("custom_app_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy: func(s *terraform.State) error {
			// Authenticators implement soft delete (deactivation), not hard delete.
			// Resources can be intentionally left in a INACTIVE state to prevent
			// accidental removal while still in use by policies or other configurations.
			// Return nil to skip destruction verification.
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName1, "type", "app"),
					resource.TestCheckResourceAttr(resourceName1, "key", "custom_app"),
					resource.TestCheckResourceAttr(resourceName1, "name", "VCRTestCustomAppAuth"),
					resource.TestCheckResourceAttr(resourceName1, "legacy_ignore_name", "false"),
					testAttributeJSON(resourceName1, "settings", `{"appInstanceId":"0oaspABCDEF12345678","userVerification":"REQUIRED"}`),
					testAttributeJSON(resourceName1, "provider_json", `{"configuration":{"fcm":{"id":"ppcrb12345678ABCDEF"}},"type":"PUSH"}`),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName1, "type", "app"),
					resource.TestCheckResourceAttr(resourceName1, "key", "custom_app"),
					resource.TestCheckResourceAttr(resourceName1, "name", "VCRTestCustomAppAuthRenamed"),
					resource.TestCheckResourceAttr(resourceName1, "legacy_ignore_name", "false"),
					testAttributeJSON(resourceName1, "settings", `{"appInstanceId":"0oaspGHIJKL12345678","userVerification":"PREFERRED"}`),
					testAttributeJSON(resourceName1, "provider_json", `{"configuration":{"fcm":{"id":"ppcrb12345678ABCDEF"}},"type":"PUSH"}`),
				),
			},
		},
	})
}
