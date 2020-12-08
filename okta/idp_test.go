package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func createDoesIdpExist(idp sdk.IdentityProvider) func(string) (bool, error) {
	return func(id string) (bool, error) {
		client := getSupplementFromMetadata(testAccProvider.Meta())
		_, response, err := client.GetIdentityProvider(context.Background(), id, idp)

		return doesResourceExist(response, err)
	}
}

func deleteTestIdps(client *testClient) error {
	var providers []*sdk.BasicIdp
	_, _, err := client.apiSupplement.ListIdentityProviders(context.Background(), &providers, &query.Params{Q: "testAcc_"})
	if err != nil {
		return err
	}

	for _, idp := range providers {
		_, err := client.apiSupplement.DeleteIdentityProvider(context.Background(), idp.Id)
		if err != nil {
			return err
		}

		if idp.Type == saml2Idp {
			_, err := client.apiSupplement.DeleteIdentityProviderSigningKey(context.Background(), idp.Protocol.Credentials.Trust.Kid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
