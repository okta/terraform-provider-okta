package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaAuthenticator_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", authenticator)
	mgr := newFixtureManager(authenticator, t.Name())
	config := mgr.GetFixtures("security_question.tf", t)
	configUpdated := mgr.GetFixtures("security_question_updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "key", "security_question"),
					resource.TestCheckResourceAttr(resourceName, "name", "Security Question"),
					testAttributeJSON(resourceName, "settings", `{"allowedFor" : "recovery"}`),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
				),
			},
		},
	})
}

// TestAccOktaAuthenticator_Issue1367_simple
// https://github.com/okta/terraform-provider-okta/issues/1367
// Google OTP is a simple example of the solution for #1367
func TestAccOktaAuthenticator_Issue1367_simple(t *testing.T) {
	config := `
resource "okta_authenticator" "google_otp" {
	name   = "Google Authenticator"
	key    = "google_otp"
	status = "INACTIVE"
}
`
	resourceName := fmt.Sprintf("%s.google_otp", authenticator)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "type", "app"),
					resource.TestCheckResourceAttr(resourceName, "key", "google_otp"),
					resource.TestCheckResourceAttr(resourceName, "name", "Google Authenticator"),
				),
			},
		},
	})
}

// TestAccOktaAuthenticator_Issue1367_provider_json
// https://github.com/okta/terraform-provider-okta/issues/1367
// Demonstrates provider input as freeform JSON
// Example from `POST /api/v1/authenticator` API docs
// https://developer.okta.com/docs/reference/api/authenticators-admin/#create-authenticator
func TestAccOktaAuthenticator_Issue1367_provider_json(t *testing.T) {
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
	resourceName := fmt.Sprintf("%s.test", authenticator)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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
