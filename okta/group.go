package okta

import "context"

func listGroupUserIds(m interface{}, id string) ([]string, error) {
	client := getOktaClientFromMetadata(m)
	arr, _, err := client.Group.ListGroupUsers(context.Background(), id, nil)
	if err != nil {
		return nil, err
	}

	userIdList := make([]string, len(arr))
	for i, user := range arr {
		userIdList[i] = user.Id
	}

	return userIdList, nil
}
