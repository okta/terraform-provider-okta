package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func createDoesIdpExist() func(string) (bool, error) {
	return func(id string) (bool, error) {
		_, response, err := getOktaClientFromMetadata(testAccProvider.Meta()).IdentityProvider.GetIdentityProvider(context.Background(), id)
		return doesResourceExist(response, err)
	}
}

func deleteTestIdps(client *testClient) error {
	providers, _, err := client.oktaClient.IdentityProvider.ListIdentityProviders(context.Background(), &query.Params{Q: "testAcc_"})
	if err != nil {
		return err
	}
	for _, idp := range providers {
		_, err := client.oktaClient.IdentityProvider.DeleteIdentityProvider(context.Background(), idp.Id)
		if err != nil {
			return err
		}

		if idp.Type == saml2Idp {
			_, err := client.oktaClient.IdentityProvider.DeleteIdentityProviderKey(context.Background(), idp.Protocol.Credentials.Trust.Kid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
