package okta

import (
	"context"
	"fmt"
	"net/url"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

func getUserTypeSchemaUrl(ctx context.Context, client *okta.Client, id string) (string, error) {
	ut, _, err := client.UserType.GetUserType(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get user type: %v", err)
	}
	return userTypeURL(ut), nil
}

func userTypeURL(ut *okta.UserType) string {
	fm, ok := ut.Links.(map[string]interface{})
	if ok {
		sm, ok := fm["schema"].(map[string]interface{})
		if ok {
			href, ok := sm["href"].(string)
			if ok {
				u, _ := url.Parse(href)
				return u.EscapedPath()
			}
		}
	}
	return ""
}
