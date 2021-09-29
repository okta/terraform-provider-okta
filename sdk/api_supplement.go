package sdk

import (
	"errors"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var errMissingAPITokenClient = errors.New("this endpoint is only available when using API token")

// APISupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type APISupplement struct {
	apiTokenClient    *okta.Client
	bearerTokenClient *okta.Client
	primary           string
}

func NewAPISupplement(apiTokenClient, bearerTokenClient *okta.Client, primary string) *APISupplement {
	return &APISupplement{
		apiTokenClient:    apiTokenClient,
		bearerTokenClient: bearerTokenClient,
		primary:           primary,
	}
}

func (m *APISupplement) RequestExecutor() *okta.RequestExecutor {
	if m.primary == "api_token_client" {
		return m.apiTokenClient.GetRequestExecutor()
	}
	return m.bearerTokenClient.GetRequestExecutor()
}
