package okta

func listGroupUserIds(m interface{}, id string) ([]string, error) {
	client := getOktaClientFromMetadata(m)
	ctx := getOktaContextFromMetadata(m)
	arr, _, err := client.Group.ListGroupUsers(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	userIdList := make([]string, len(arr))
	for i, user := range arr {
		userIdList[i] = user.Id
	}

	return userIdList, nil
}
