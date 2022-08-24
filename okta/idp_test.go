package okta

import (
	"context"
)

func createDoesIdpExist() func(string) (bool, error) {
	return func(id string) (bool, error) {
		_, response, err := getOktaClientFromMetadata(testAccProvider.Meta()).IdentityProvider.GetIdentityProvider(context.Background(), id)
		return doesResourceExist(response, err)
	}
}
