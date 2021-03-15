package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAdminRoleTargets(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(adminRoleTargets)
	basic := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceAppName := fmt.Sprintf("%s.test_app", adminRoleTargets)
	resourceGroupName := fmt.Sprintf("%s.test_group", adminRoleTargets)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(idpOidc, doesTargetExists()),
		Steps: []resource.TestStep{
			{
				Config: basic,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceAppName, doesTargetExists()),
					ensureResourceExists(resourceGroupName, doesTargetExists()),
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
					ensureResourceExists(resourceAppName, doesTargetExists()),
					ensureResourceExists(resourceGroupName, doesTargetExists()),
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

func doesTargetExists() func(string) (bool, error) {
	return func(id string) (bool, error) {
		parts := strings.Split(id, "/")
		roles, _, err := getOktaClientFromMetadata(testAccProvider.Meta()).User.ListAssignedRolesForUser(context.Background(), parts[0], nil)
		if err != nil {
			return false, fmt.Errorf("failed to get list of roles associated with the user: %v", err)
		}
		for i := range roles {
			if roles[i].Type != parts[1] {
				continue
			}
			apps, _, err := getOktaClientFromMetadata(testAccProvider.Meta()).User.
				ListApplicationTargetsForApplicationAdministratorRoleForUser(
					context.Background(), parts[0], roles[i].Id, nil)
			if err != nil {
				return false, fmt.Errorf("failed to read app targets: %v", err)
			}
			if len(apps) > 0 {
				return true, nil
			}
			groups, _, err := getOktaClientFromMetadata(testAccProvider.Meta()).User.
				ListGroupTargetsForRole(context.Background(), parts[0], roles[i].Id, nil)
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
}
