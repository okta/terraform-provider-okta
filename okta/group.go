package okta

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func listGroupUserIDs(ctx context.Context, m interface{}, id string) ([]string, error) {
	var resUsers []string
	users, resp, err := getOktaClientFromMetadata(m).Group.ListGroupUsers(ctx, id, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return nil, err
	}
	for {
		for _, user := range users {
			resUsers = append(resUsers, user.Id)
		}
		if resp.HasNextPage() {
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

func listGroups(ctx context.Context, client *okta.Client, qp *query.Params) ([]*okta.Group, error) {
	var resGroups []*okta.Group
	groups, resp, err := client.Group.ListGroups(ctx, qp)
	if err != nil {
		return nil, err
	}
	for {
		resGroups = append(resGroups, groups...)
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &groups)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			break
		}
	}
	return resGroups, nil
}

// Group Primary Key Operations
func addGroupMembers(ctx context.Context, client *okta.Client, groupId string, users []string) error {
	for _, user := range users {
		resp, err := client.Group.AddUserToGroup(ctx, groupId, user)
		exists, err := doesResourceExist(resp, err)
		if err != nil {
			return fmt.Errorf("failed to add user (%s) to group (%s): %v", user, groupId, err)
		}
		if !exists {
			return fmt.Errorf("targeted object does not exist: %s", err)
		}
	}
	return nil
}

func removeGroupMembers(ctx context.Context, client *okta.Client, groupId string, users []string) error {
	for _, user := range users {
		resp, err := client.Group.RemoveUserFromGroup(ctx, groupId, user)
		err = suppressErrorOn404(resp, err)
		if err != nil {
			return fmt.Errorf("failed to remove user (%s) from group (%s): %v", user, groupId, err)
		}
	}
	return nil
}

//User Primary Key Operatios
func addUserToGroups(ctx context.Context, client *okta.Client, userId string, groups []string) error {
	for _, group := range groups {
		resp, err := client.Group.AddUserToGroup(ctx, group, userId)
		exists, err := doesResourceExist(resp, err)
		if err != nil {
			return fmt.Errorf("failed to add user (%s) to group (%s): %v", userId, group, err)
		}
		if !exists {
			return fmt.Errorf("targeted object does not exist: %s", err)
		}
	}
	return nil
}

func removeUserFromGroups(ctx context.Context, client *okta.Client, userId string, groups []string) error {
	for _, group := range groups {
		resp, err := client.Group.RemoveUserFromGroup(ctx, group, userId)
		err = suppressErrorOn404(resp, err)
		if err != nil {
			return fmt.Errorf("failed to remove user (%s) from group (%s): %v", userId, group, err)
		}
	}
	return nil
}
