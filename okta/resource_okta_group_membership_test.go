package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaGroupMembership_crud(t *testing.T) {
	mgr := newFixtureManager(groupMembership, t.Name())
	config := mgr.GetFixtures("okta_group_membership.tf", t)
	updatedConfig := mgr.GetFixtures("okta_group_membership_updated.tf", t)
	removedConfig := mgr.GetFixtures("okta_group_membership_removed.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(groupMembership, checkMembershipState),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: updatedConfig,
			},
			{
				Config: removedConfig,
			},
		},
	})
}

func checkMembershipState(id string) (bool, error) {
	ids := strings.Split(id, "+")
	groupId := ids[0]
	userId := ids[1]
	client := getOktaClientFromMetadata(testAccProvider.Meta())
	state, err := checkIfUserInGroup(context.Background(), client, groupId, userId)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("Resource not found: %s (UserGroup)", groupId)) {
			return state, nil
		}
	}
	return state, err
}
