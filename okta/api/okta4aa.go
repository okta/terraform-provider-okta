package api

import (
	okta4AI "github.com/okta/okta-for-ai-agents-golang/okta4aa"
)

var _ Okta4AIClient = &okta4AIAPIClient{}

type okta4AIAPIClient struct {
	okta4AISDKClient *okta4AI.APIClient
}

type Okta4AIClient interface {
	Okta4AISDKClient() *okta4AI.APIClient
}

func newOkta4AISDKClient(c *OktaAPIConfig) (*okta4AI.APIClient, error) {
	config, _, err := getO4AIClientConfig(c)
	if err != nil {
		return nil, err
	}
	return okta4AI.NewAPIClient(config), nil
}

func (c *okta4AIAPIClient) Okta4AISDKClient() *okta4AI.APIClient {
	return c.okta4AISDKClient
}

func NewOkta4AIAPIClient(c *OktaAPIConfig) (Okta4AIClient, error) {
	client, err := newOkta4AISDKClient(c)
	if err != nil {
		return nil, err
	}
	return &okta4AIAPIClient{okta4AISDKClient: client}, nil
}
