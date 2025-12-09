package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

// Note: at least one factor (e.g. `okta_otp`) should be enabled before running this test.
func TestAccResourceOktaMfaPolicy_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyMfa, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyMfa)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkPolicyDestroy(resources.OktaIDaaSPolicyMfa),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "google_otp.enroll", "REQUIRED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_new"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Updated"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "google_otp.enroll", "OPTIONAL"),
				),
			},
		},
	})
}

// TestAccResourceOktaMfaPolicy_PR_1210 deals with testing
// https://github.com/okta/terraform-provider-okta/pull/1210
func TestAccResourceOktaMfaPolicy_PR_1210(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyMfa, t.Name())
	config := `
data "okta_group" "all" {
  name = "Everyone"
}
resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  status = "ACTIVE"
  description = "Terraform Acceptance Test MFA Policy"
  priority = 1
  is_oie  = true

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_email = {
    enroll = "NOT_ALLOWED"
  }

  fido_webauthn = {
    enroll = "REQUIRED"
  }

  groups_included = [data.okta_group.all.id]
}
	`
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyMfa)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkPolicyDestroy(resources.OktaIDaaSPolicyMfa),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "NOT_ALLOWED"),
					resource.TestCheckResourceAttr(resourceName, "fido_webauthn.enroll", "REQUIRED"),
				),
			},
		},
	})
}

// TestAccResourceOktaMfaPolicy_Issue_1176 deals with testing
// https://github.com/okta/terraform-provider-okta/issues/1176
// Which is similar to PRs 1427/1210
func TestAccResourceOktaMfaPolicy_Issue_1176(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyMfa, t.Name())
	config := `
data "okta_group" "all" {
  name = "Everyone"
}
resource "okta_policy_mfa" "test" {
    name        = "testAcc_replace_with_uuid"
    status      = "ACTIVE"
    description = "Terraform Acceptance Test MFA Policy"
    is_oie      = true
    okta_otp = {
      enroll = "OPTIONAL"
    }
    #phone_number = {
    #  enroll = "OPTIONAL"
    #}
    okta_password = {
      enroll = "REQUIRED"
    }
    okta_email = {
      enroll = "OPTIONAL"
    }

    groups_included = [data.okta_group.all.id]
}
	`
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyMfa)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkPolicyDestroy(resources.OktaIDaaSPolicyMfa),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy"),
					resource.TestCheckResourceAttr(resourceName, "okta_otp.enroll", "OPTIONAL"),
					// phone authentictor needs to be enabled on this org to make this acc test pass all the time
					// resource.TestCheckResourceAttr(resourceName, "phone_number.enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "OPTIONAL"),
				),
			},
		},
	})
}

func TestAccResourceOktaMfaPolicy_Issue_2139_yubikey_token(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyMfa, t.Name())
	config := `
data "okta_group" "all" {
  name = "Everyone"
}
resource "okta_policy_mfa" "test" {
    name        = "testAcc_replace_with_uuid"
    description = "Terraform Acceptance Test MFA Policy Yubikey Token"
    status      = "ACTIVE"
    is_oie      = true
    groups_included = [data.okta_group.all.id]
    okta_password = {
      enroll = "REQUIRED"
    }
    yubikey_token = {
      enroll = "%s"
    }
}
	`
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyMfa)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkPolicyDestroy(resources.OktaIDaaSPolicyMfa),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(config, "OPTIONAL")),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Yubikey Token"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "yubikey_token.enroll", "OPTIONAL"),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(config, "REQUIRED")),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Yubikey Token"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "yubikey_token.enroll", "REQUIRED"),
				),
			},
		},
	})
}

func TestAccResourceOktaMfaPolicy_custom_otp(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyMfa, t.Name())
	config := `
data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
    name        = "testAcc_replace_with_uuid"
    description = "Terraform Acceptance Test MFA Policy Custom OTP"
    status      = "ACTIVE"
    is_oie      = true
    groups_included = [data.okta_group.all.id]
    okta_password = {
      enroll = "REQUIRED"
    }
    custom_otp = {
      enroll = "%s"
    }
}
	`
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyMfa)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkPolicyDestroy(resources.OktaIDaaSPolicyMfa),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(config, "NOT_ALLOWED")),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Custom OTP"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "custom_otp.enroll", "NOT_ALLOWED"),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(config, "OPTIONAL")),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Custom OTP"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "custom_otp.enroll", "OPTIONAL"),
				),
			},
		},
	})
}
