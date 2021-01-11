package okta

import "context"

func listGroupUserIDs(ctx context.Context, m interface{}, id string) ([]string, error) {
	arr, _, err := getOktaClientFromMetadata(m).Group.ListGroupUsers(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	userIDList := make([]string, len(arr))
	for i, user := range arr {
		userIDList[i] = user.Id
	}

	return userIDList, nil
}
