package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaEmailDomainVerification_alreadyVerified(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSEmailDomainVerification, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSEmailDomainVerification)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email_domain_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: importStateIdForEmailDomainVerification(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func importStateIdForEmailDomainVerification(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["email_domain_id"], nil
	}
}
