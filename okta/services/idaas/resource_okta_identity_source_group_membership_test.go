package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaIdentitySourceGroupMembership_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSIdentitySourceGroupMembership, t.Name())
	config := mgr.GetFixtures("resource.tf", t)
	resourceName := fmt.Sprintf("%s.example", resources.OktaIDaaSIdentitySourceGroupMembership)

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
					resource.TestCheckResourceAttr(resourceName, "group_or_external_id", "GRPEXT123456TESTGROUP1"),
					resource.TestCheckResourceAttr(resourceName, "member_external_id", "USEREXT123456TESTUSER1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       importStateIdForIdentitySourceGroupMembership(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"member_external_id"},
			},
		},
	})
}

func importStateIdForIdentitySourceGroupMembership(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		identitySourceId := rs.Primary.Attributes["identity_source_id"]
		groupOrExternalId := rs.Primary.Attributes["group_or_external_id"]
		id := rs.Primary.ID
		return fmt.Sprintf("%s/%s/%s", identitySourceId, groupOrExternalId, id), nil
	}
}
