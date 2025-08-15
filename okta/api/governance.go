package api

import (
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
)

type governanceAPIClient struct {
	oktaIGSDKClient *governance.IGAPIClient
}

type OktaGovernanceClient interface {
	OktaGovernanceSDKClient() *governance.IGAPIClient
}

func oktaGovernanceSDKClient(c *OktaAPIConfig) (client *governance.IGAPIClient, err error) {
	err, config, _, _ := getV5ClientConfig(c, err)

	client = governance.NewAPIClient(config)
	return client, nil
}

func (c *governanceAPIClient) OktaGovernanceSDKClient() *governance.IGAPIClient {
	return c.oktaIGSDKClient
}

func NewOktaGovernanceAPIClient(c *OktaAPIConfig) (client OktaGovernanceClient, err error) {
	governanceClient, err := oktaGovernanceSDKClient(c)
	if err != nil {
		return
	}

	client = &governanceAPIClient{
		oktaIGSDKClient: governanceClient,
	}

	return
}
