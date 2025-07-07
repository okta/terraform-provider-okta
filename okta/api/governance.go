package api

import (
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
)

type governanceAPIClient struct {
	oktaIGSDKClientV5 *oktaInternalGovernance.IGAPIClient
}

type OktaGovernanceClient interface {
	OktaIGSDKClientV5() *oktaInternalGovernance.IGAPIClient
}

func oktaV5IGSDKClient(c *OktaAPIConfig) (client *oktaInternalGovernance.IGAPIClient, err error) {
	err, config, _, _ := getV5ClientConfig(c, err)

	client = oktaInternalGovernance.NewAPIClient(config)
	return client, nil
}

func (c *governanceAPIClient) OktaIGSDKClientV5() *oktaInternalGovernance.IGAPIClient {
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
