package api

import (
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
)

type governanceAPIClient struct {
	oktaIGSDKClientV5 *governance.IGAPIClient
}

type OktaGovernanceClient interface {
	OktaIGSDKClient() *governance.IGAPIClient
}

func oktaIGSDKClient(c *OktaAPIConfig) (client *governance.IGAPIClient, err error) {
	err, config, _, _ := getV5ClientConfig(c, err)

	client = governance.NewAPIClient(config)
	return client, nil
}

func (c *governanceAPIClient) OktaIGSDKClient() *governance.IGAPIClient {
	return c.oktaIGSDKClientV5
}

func NewOktaGovernanceAPIClient(c *OktaAPIConfig) (client OktaGovernanceClient, err error) {
	v5IGClient, err := oktaIGSDKClient(c)
	if err != nil {
		return
	}

	client = &governanceAPIClient{
		oktaIGSDKClientV5: v5IGClient,
	}

	return
}
