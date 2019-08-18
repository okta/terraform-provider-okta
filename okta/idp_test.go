package okta

import (
	"github.com/articulate/terraform-provider-okta/sdk"
	"github.com/okta/okta-sdk-golang/okta/query"
)

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
