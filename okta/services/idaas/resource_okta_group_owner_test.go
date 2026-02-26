package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccResourceOktaGroupOwner_crud(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owner", t.Name())
	config := mgr.GetFixtures("resource.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("okta_user.test", "first_name", "TestAcc"),
						resource.TestCheckResourceAttr("okta_user.test", "last_name", "Smith"),
						resource.TestCheckResourceAttr("okta_group_owner.test", "type", "USER"),
					),
				},
				{
					ResourceName:      "okta_group_owner.test",
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						groupID := s.RootModule().Resources["okta_group.test"].Primary.Attributes["id"]
						groupOwnerID := s.RootModule().Resources["okta_group_owner.test"].Primary.Attributes["id"]
						return fmt.Sprintf("%s/%s", groupID, groupOwnerID), nil
					},
				},
			},
		})
}
