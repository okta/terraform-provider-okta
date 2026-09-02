package api

import (
	okta4AI "github.com/okta4AI/okta4AI-for-ai-sdk-golang/v1"
)

var _ Okta4AIClient = &okta4AIAPIClient{}

type okta4AIAPIClient struct {
	okta4AISDKClient *okta4AI.Okta4AIAPIClient
}

type Okta4AIClient interface {
	Okta4AISDKClient() *okta4AI.Okta4AIAPIClient
}

func newOkta4AISDKClient(c *OktaAPIConfig) (*okta4AI.Okta4AIAPIClient, error) {
	config, _, err := getO4AIClientConfig(c)
	if err != nil {
		return nil, err
	}
	return okta4AI.NewAPIClient(config), nil
}

func (c *okta4AIAPIClient) Okta4AISDKClient() *okta4AI.Okta4AIAPIClient {
	return c.okta4AISDKClient
}

func NewOkta4AIAPIClient(c *OktaAPIConfig) (Okta4AIClient, error) {
	client, err := newOkta4AISDKClient(c)
	if err != nil {
		return nil, err
	}
	return &okta4AIAPIClient{okta4AISDKClient: client}, nil
}
