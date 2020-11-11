package okta

import "context"

func listGroupUserIDs(m interface{}, id string) ([]string, error) {
	client := getOktaClientFromMetadata(m)
	arr, _, err := client.Group.ListGroupUsers(context.Background(), id, nil)
	if err != nil {
		return nil, err
	}

	userIDList := make([]string, len(arr))
	for i, user := range arr {
		userIDList[i] = user.Id
	}

	return userIDList, nil
}
