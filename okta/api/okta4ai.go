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

// TODO(okta4ai): wire the real configuration. getV5ClientConfig returns the v5 SDK's
// *Configuration, but okta4AI.NewAPIClient takes the okta4AI package's own *Configuration
// (built via okta4AI.NewConfiguration(...ConfigSetter), which returns two values). Those two
// types are unrelated, so the previous call here did not compile. Stubbed so that the rest of
// the provider type-checks; the okta4AI client is non-functional until this is implemented.
func newOkta4AISDKClient(c *OktaAPIConfig) (*okta4AI.APIClient, error) {
	_ = c
	return nil, nil
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
