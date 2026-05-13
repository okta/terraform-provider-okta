package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaIdentitySourceGroup_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSIdentitySourceGroup, t.Name())
	config := mgr.GetFixtures("resource.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSIdentitySourceGroup)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "identity_source_id", "0oaxc95befZNgrJl71d7"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "GRPEXT123456TESTGROUP1"),
					resource.TestCheckResourceAttr(resourceName, "profile.display_name", "Test Engineering Group"),
					resource.TestCheckResourceAttr(resourceName, "profile.description", "A test group for identity source integration"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: importStateIdForIdentitySourceGroup(resourceName),
				ImportStateVerify: true,
				// profile is not returned by the GET endpoint — Read cannot restore it
				ImportStateVerifyIgnore: []string{"profile"},
			},
		},
	})
}

func importStateIdForIdentitySourceGroup(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		identitySourceId := rs.Primary.Attributes["identity_source_id"]
		id := rs.Primary.ID
		return fmt.Sprintf("%s/%s", identitySourceId, id), nil
	}
}
