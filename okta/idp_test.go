package okta

import (
	"context"
)

func createDoesIdpExist(id string) (bool, error) {
	client := sdkV2ClientForTest()
	_, response, err := client.IdentityProvider.GetIdentityProvider(context.Background(), id)
	return doesResourceExist(response, err)
}
