package idaas_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAdminRoleTargets_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAdminRoleTargets, t.Name())
	basic := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceAppName := fmt.Sprintf("%s.test_app", resources.OktaIDaaSAdminRoleTargets)
	resourceGroupName := fmt.Sprintf("%s.test_group", resources.OktaIDaaSAdminRoleTargets)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSIdpOidc, doesTargetExists),
		Steps: []resource.TestStep{
			{
				Config: basic,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceAppName, doesTargetExists),
					ensureResourceExists(resourceGroupName, doesTargetExists),
					resource.TestCheckResourceAttrSet(resourceAppName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceAppName, "role_id"),
					resource.TestCheckResourceAttr(resourceAppName, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceAppName, "apps.#", "1"),
					resource.TestCheckResourceAttrSet(resourceGroupName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceGroupName, "role_id"),
					resource.TestCheckResourceAttr(resourceGroupName, "role_type", "GROUP_MEMBERSHIP_ADMIN"),
					resource.TestCheckResourceAttr(resourceGroupName, "groups.#", "1"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceAppName, doesTargetExists),
					ensureResourceExists(resourceGroupName, doesTargetExists),
					resource.TestCheckResourceAttrSet(resourceAppName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceAppName, "role_id"),
					resource.TestCheckResourceAttr(resourceAppName, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceAppName, "apps.#", "2"),
					resource.TestCheckResourceAttrSet(resourceGroupName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceGroupName, "role_id"),
					resource.TestCheckResourceAttr(resourceGroupName, "role_type", "GROUP_MEMBERSHIP_ADMIN"),
					resource.TestCheckResourceAttr(resourceGroupName, "groups.#", "2"),
				),
			},
		},
	})
}

func doesTargetExists(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	parts := strings.Split(id, "/")
	roles, _, err := client.User.ListAssignedRolesForUser(context.Background(), parts[0], nil)
	if err != nil {
		return false, fmt.Errorf("failed to get list of roles associated with the user: %v", err)
	}
	for i := range roles {
		if roles[i].Type != parts[1] {
			continue
		}
		apps, _, err := client.User.ListApplicationTargetsForApplicationAdministratorRoleForUser(
			context.Background(), parts[0], roles[i].Id, nil)
		if err != nil {
			return false, fmt.Errorf("failed to read app targets: %v", err)
		}
		if len(apps) > 0 {
			return true, nil
		}
		groups, _, err := client.User.ListGroupTargetsForRole(context.Background(), parts[0], roles[i].Id, nil)
		if err != nil {
			return false, fmt.Errorf("failed to read group targets: %v", err)
		}
		if len(groups) > 0 {
			return true, nil
		}
		break
	}
	return false, nil
}
