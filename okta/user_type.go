package okta

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/okta/terraform-provider-okta/sdk"
)

func getUserTypeSchemaID(ctx context.Context, client *sdk.Client, id string) (string, error) {
	if id == "default" {
		return "default", nil
	}
	ut, _, err := client.UserType.GetUserType(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get user type: %v", err)
	}
	return userTypeSchemaID(ut), nil
}

func userTypeSchemaID(ut *sdk.UserType) string {
	fm, ok := ut.Links.(map[string]interface{})
	if ok {
		sm, ok := fm["schema"].(map[string]interface{})
		if ok {
			href, ok := sm["href"].(string)
			if ok {
				u, _ := url.Parse(href)
				return strings.TrimPrefix(u.EscapedPath(), "/api/v1/meta/schemas/user/")
			}
		}
	}
	return ""
}
