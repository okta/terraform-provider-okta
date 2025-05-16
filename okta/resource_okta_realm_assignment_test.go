package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaRealmAssignment_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", realmAssignment)
	mgr := newFixtureManager("resources", realmAssignment, t.Name())
	config := mgr.GetFixtures("okta_realm_assignment.tf", t)
	updatedConfig := mgr.GetFixtures("okta_realm_assignment_updated.tf", t)
	buildResourceName(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		CheckDestroy:             checkResourceDestroy(realmAssignment, doesRealmAssignmentExist),
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
	client := sdkV5ClientForTest()
	_, response, err := client.RealmAssignmentAPI.GetRealmAssignment(context.Background(), id).Execute()
	return doesResourceExistV5(response, err)
}
