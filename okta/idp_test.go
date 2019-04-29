package okta

import (
	"github.com/okta/okta-sdk-golang/okta/query"
)

func deleteTestIdps(client *testClient) error {
	providers := []*BasicIdp{}
	_, _, err := client.apiSupplement.ListIdentityProviders(providers, &query.Params{Q: "testAcc_"})
	if err != nil {
		return err
	}

	for _, idp := range providers {
		_, err := client.apiSupplement.DeleteIdentityProvider(idp.Id)
		if err != nil {
			return err
		}
	}

	return nil
}
