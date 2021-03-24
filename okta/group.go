package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func listGroupUserIDs(ctx context.Context, m interface{}, id string) ([]string, error) {
	var resUsers []string
	users, resp, err := getOktaClientFromMetadata(m).Group.ListGroupUsers(ctx, id, &query.Params{Limit: 200})
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
