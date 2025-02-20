package idaas_test

import (
	"context"

	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func createDoesIdpExist(id string) (bool, error) {
	client := provider.SdkV2ClientForTest()
	_, response, err := client.IdentityProvider.GetIdentityProvider(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}
