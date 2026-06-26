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

func TestDomainValidationString(t *testing.T) {
	tests := []struct {
		element  string
		expected bool
	}{
		{"VERIFIED", true},
		{"COMPLETED", true},
		{"NOT_STARTED", false},
		{"IN_PROGRESS", false},
		{"verified", false},
		{"completed", false},
		{"invalid", false},
	}

	for _, test := range tests {
		actual := idaas.IsDomainValidated(test.element)

		if actual != test.expected {
			t.Errorf("domain validation status failed for status = \"%s\" - Expected: %t, Actual: %t", test.element, test.expected, actual)
		}
	}
}

func TestAccResourceOktaDomainVerification_import(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSDomainVerification, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSDomainVerification)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: importStateIdForDomainVerification(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func importStateIdForDomainVerification(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["domain_id"], nil
	}
}
