package okta

import (
	"github.com/okta/okta-sdk-golang/okta/query"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func createDoesIdpExist(idp sdk.IdentityProvider) func(string) (bool, error) {
	return func(id string) (bool, error) {
		client := getSupplementFromMetadata(testAccProvider.Meta())
		_, response, err := client.GetIdentityProvider(id, idp)

		return doesResourceExist(response, err)
	}
}

func deleteTestIdps(client *testClient) error {
	providers := []*sdk.BasicIdp{}
	_, _, err := client.apiSupplement.ListIdentityProviders(&providers, &query.Params{Q: "testAcc_"})
	if err != nil {
		return err
	}

	for _, idp := range providers {
		_, err := client.apiSupplement.DeleteIdentityProvider(idp.Id)
		if err != nil {
			return err
		}

		if idp.Type == saml2Idp {
			_, err := client.apiSupplement.DeleteIdentityProviderSigningKey(idp.Protocol.Credentials.Trust.Kid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
