package api

import (
	"github.com/okta/okta-governance-sdk-golang/governance"
)

type governanceAPIClient struct {
	oktaGovernanceSDKClient *governance.OktaGovernanceAPIClient
}

type OktaGovernanceClient interface {
	OktaGovernanceSDKClient() *governance.OktaGovernanceAPIClient
}

func oktaGovernanceSDKClient(c *OktaAPIConfig) (client *governance.OktaGovernanceAPIClient, err error) {
	config, _, _ := getV5ClientConfig(c)
	client = governance.NewAPIClient(config)
	return client, nil
}

func (c *governanceAPIClient) OktaGovernanceSDKClient() *governance.OktaGovernanceAPIClient {
	return c.oktaGovernanceSDKClient
}

func NewOktaGovernanceAPIClient(c *OktaAPIConfig) (client OktaGovernanceClient, err error) {
	governanceClient, err := oktaGovernanceSDKClient(c)
	if err != nil {
		return
	}

	client = &governanceAPIClient{
		oktaGovernanceSDKClient: governanceClient,
	}

	return
}
