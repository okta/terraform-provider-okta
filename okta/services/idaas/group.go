package idaas

import (
	"context"
	"fmt"

	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func listGroupUsers(ctx context.Context, m interface{}, id string) ([]*sdk.User, error) {
	var resUsers []*sdk.User
	users, resp, err := getOktaClientFromMetadata(m).Group.ListGroupUsers(ctx, id, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return nil, err
	}
	for {
		resUsers = append(resUsers, users...)
		if resp.HasNextPage() {
			users = []*sdk.User{}
			resp, err = resp.Next(ctx, &users)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			break
		}
	}
	return resUsers, nil
}

func listGroupUserIDs(ctx context.Context, m interface{}, id string) ([]string, error) {
	var resUsers []string
	users, err := listGroupUsers(ctx, m, id)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		resUsers = append(resUsers, user.Id)
	}
	return resUsers, nil
}

// Group Primary Key Operations (Use when # groups < # users in operations)
func addGroupMembers(ctx context.Context, client *v6okta.APIClient, groupId string, users []string) error {
	for _, user := range users {
		resp, err := client.GroupAPI.AssignUserToGroup(ctx, groupId, user).Execute()
		if err != nil {
			return fmt.Errorf("failed to add user (%s) to group (%s): %w", user, groupId, err)
		}
		exists, err := utils.DoesResourceExistV6(resp, err)
		if !exists {
			return fmt.Errorf("targeted object does not exist: %s", err)
		}
	}
	return nil
}

func removeGroupMembers(ctx context.Context, client *v6okta.APIClient, groupId string, users []string) error {
	for _, user := range users {
		resp, err := client.GroupAPI.UnassignUserFromGroup(ctx, groupId, user).Execute()
		err = utils.SuppressErrorOn404_V6(resp, err)
		if err != nil {
			return fmt.Errorf("failed to remove user (%s) from group (%s): %v", user, groupId, err)
		}
	}
	return nil
}

// User Primary Key Operations (use when # users < # groups in operations)
func addUserToGroups(ctx context.Context, client *v6okta.APIClient, userId string, groups []string) error {
	for _, group := range groups {
		resp, err := client.GroupAPI.AssignUserToGroup(ctx, group, userId).Execute()
		exists, err := utils.DoesResourceExistV6(resp, err)
		if err != nil {
			return fmt.Errorf("failed to add user (%s) to group (%s): %v", userId, group, err)
		}
		if !exists {
			return fmt.Errorf("targeted object does not exist: %s", err)
		}
	}
	return nil
}

func removeUserFromGroups(ctx context.Context, client *v6okta.APIClient, userId string, groups []string) error {
	for _, group := range groups {
		resp, err := client.GroupAPI.UnassignUserFromGroup(ctx, group, userId).Execute()
		err = utils.SuppressErrorOn404_V6(resp, err)
		if err != nil {
			return fmt.Errorf("failed to remove user (%s) from group (%s): %v", userId, group, err)
		}
	}
	return nil
}
