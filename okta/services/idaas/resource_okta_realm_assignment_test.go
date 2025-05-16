package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaRealmAssignment_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSRealmAssignment)
	mgr := newFixtureManager("resources", resources.OktaIDaaSRealmAssignment, t.Name())

	config := mgr.GetFixtures("okta_realm_assignment.tf", t)
	updatedConfig := mgr.GetFixtures("okta_realm_assignment_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSRealmAssignment, doesRealmAssignmentExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "realm_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "TestAcc Example Realm Assignment"),
					resource.TestCheckResourceAttr(resourceName, "priority", "55"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "condition_expression", "user.profile.login.contains(\"@acctest.com\")"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "TestAcc Example Realm Assignment Updated"),
					resource.TestCheckResourceAttr(resourceName, "priority", "111"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "condition_expression", "user.profile.login.contains(\"@acctestupdated.com\")"),
				),
			},
		},
	})
}

func doesRealmAssignmentExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV5()
	_, response, err := client.RealmAssignmentAPI.GetRealmAssignment(context.Background(), id).Execute()
	return utils.DoesResourceExistV5(response, err)
}
