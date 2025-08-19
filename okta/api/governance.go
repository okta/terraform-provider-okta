package api

import (
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
)

type governanceAPIClient struct {
	oktaGovernanceSDKClient *governance.OktaGovernanceAPIClient
}

type OktaGovernanceClient interface {
	OktaGovernanceSDKClient() *governance.OktaGovernanceAPIClient
}

func oktaV5IGSDKClient(c *OktaAPIConfig) (client *governance.OktaGovernanceAPIClient, err error) {
	err, config, _, _ := getV5ClientConfig(c, err)

	client = governance.NewAPIClient(config)
	return client, nil
}

func (c *governanceAPIClient) OktaGovernanceSDKClient() *governance.OktaGovernanceAPIClient {
	return c.oktaGovernanceSDKClient
}

func NewOktaGovernanceAPIClient(c *OktaAPIConfig) (client OktaGovernanceClient, err error) {
	v5IGClient, err := oktaV5IGSDKClient(c)
	if err != nil {
		return
	}

	client = &governanceAPIClient{
		oktaGovernanceSDKClient: v5IGClient,
	}

	return
}
