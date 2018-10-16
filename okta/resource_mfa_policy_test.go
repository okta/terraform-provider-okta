package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func deleteMfaPolicies(client *testClient) error {
	return deletePolicyByType(mfaPolicyType, client)
}

func TestAccOktaMfaPolicy(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaMfaPolicy(ri)
	updatedConfig := testOktaMfaPolicyUpdated(ri)
	resourceName := buildResourceFQN(mfaPolicy, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createPolicyCheckDestroy(mfaPolicy),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy"),
					resource.TestCheckResourceAttr(resourceName, "google_otp_enroll", "REQUIRED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Updated"),
					resource.TestCheckResourceAttr(resourceName, "duo_enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "fido_u2f_enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "fido_webauthn_enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "google_otp_enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "okta_otp_enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "okta_sms_enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "rsa_token_enroll", "OPTIONAL"),
				),
			},
		},
	})
}

func testOktaMfaPolicy(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "okta_factor" "google" {
	provider_id = "google_otp" 
}
resource "%s" "%s" {
  name        		= "%s"
  status      		= "ACTIVE"
  description 		= "Terraform Acceptance Test MFA Policy"
  google_otp_enroll = "REQUIRED"
}
`, mfaPolicy, name, name)
}

func testOktaMfaPolicyUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_everyone_group" "everyone-%d" {}
resource "okta_factor" "google" {
	provider_id = "google_otp" 
}
resource "okta_factor" "sms" {
	provider_id = "okta_sms" 
}
resource "okta_factor" "otp" {
	provider_id = "okta_otp" 
}
resource "okta_factor" "duo" {
	provider_id = "duo" 
}
resource "okta_factor" "fido_u2f" {
	provider_id = "fido_u2f" 
}
resource "okta_factor" "fido_webauthn" {
	provider_id = "fido_webauthn" 
}
resource "okta_factor" "rsa_token" {
	provider_id = "rsa_token" 
}

resource "%s" "%s" {
	name        = "%s"
	status      = "INACTIVE"
	description = "Terraform Acceptance Test MFA Policy Updated"
	groups_included = [ "${data.okta_everyone_group.everyone-%d.id}" ]
	duo_enroll 				= "OPTIONAL"
	fido_u2f_enroll 		= "OPTIONAL"
	fido_webauthn_enroll 	= "OPTIONAL"
	google_otp_enroll	 	= "OPTIONAL"
	okta_otp_enroll 		= "OPTIONAL"
	okta_sms_enroll 		= "OPTIONAL"
	rsa_token_enroll	 	= "OPTIONAL"
}
`, rInt, mfaPolicy, name, name, rInt)
}
