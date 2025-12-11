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
	config := mgr.GetFixtures("pr_1210.tf", t)
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
	config := mgr.GetFixtures("issue_1176.tf", t)
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
	config := mgr.GetFixtures("issue_2139_yubikey_token.tf", t)
	configRequired := mgr.GetFixtures("issue_2139_yubikey_token_required.tf", t)
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
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Yubikey Token"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "yubikey_token.enroll", "OPTIONAL"),
				),
			},
			{
				Config: configRequired,
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

func TestAccResourceOktaMfaPolicy_custom_app(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyMfa, t.Name())
	config := mgr.GetFixtures("custom_app.tf", t)
	configModified := mgr.GetFixtures("custom_app_modified.tf", t)
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
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy with Specific Custom Apps"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "security_question.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.0.id", "aut123456789abcdef"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.0.enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.1.id", "aut123456789ghijkl"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.1.enroll", "OPTIONAL"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy with Specific Custom Apps"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "security_question.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.0.id", "aut123456789abcdef"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.0.enroll", "NOT_ALLOWED"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.1.id", "aut123456789ghijkl"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.1.enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.2.id", "aut123456789mnopqr"),
					resource.TestCheckResourceAttr(resourceName, "custom_app.2.enroll", "OPTIONAL"),
				),
			},
		},
	})
}
