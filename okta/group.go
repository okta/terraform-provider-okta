package okta

import (
	"context"

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
