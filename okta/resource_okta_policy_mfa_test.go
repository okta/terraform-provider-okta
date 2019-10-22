package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func deleteMfaPolicies(client *testClient) error {
	return deletePolicyByType(mfaPolicyType, client)
}

func TestAccOktaMfaPolicy_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyMfa)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyMfa)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createPolicyCheckDestroy(policyMfa),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy"),
					resource.TestCheckResourceAttr(resourceName, "okta_otp.enroll", "REQUIRED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Updated"),
					resource.TestCheckResourceAttr(resourceName, "fido_u2f.enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "okta_otp.enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "okta_sms.enroll", "OPTIONAL"),
				),
			},
		},
	})
}
