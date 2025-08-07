package api

import (
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
)

type governanceAPIClient struct {
	oktaIGSDKClientV5 *governance.IGAPIClient
}

type OktaGovernanceClient interface {
	OktaIGSDKClientV5() *governance.IGAPIClient
}

func oktaV5IGSDKClient(c *OktaAPIConfig) (client *governance.IGAPIClient, err error) {
	err, config, _, _ := getV5ClientConfig(c, err)

	client = governance.NewAPIClient(config)
	return client, nil
}

func (c *governanceAPIClient) OktaIGSDKClientV5() *governance.IGAPIClient {
	return c.oktaIGSDKClientV5
}

func NewOktaGovernanceAPIClient(c *OktaAPIConfig) (client OktaGovernanceClient, err error) {
	v5IGClient, err := oktaV5IGSDKClient(c)
	if err != nil {
		return
	}

	client = &governanceAPIClient{
		oktaIGSDKClientV5: v5IGClient,
	}

	return
}
