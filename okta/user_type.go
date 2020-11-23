package okta

import (
	"context"
	"net/url"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

func getUserTypeSchemaUrl(m interface{}, id string) (string, error) {
	ut, _, err := getOktaClientFromMetadata(m).UserType.GetUserType(context.Background(), id)
	if err != nil {
		return "", err
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
