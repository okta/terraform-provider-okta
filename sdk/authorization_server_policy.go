package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

func (m *ApiSupplement) ActivateAuthorizationServerPolicy(ctx context.Context, authServerID, policyID string) (*okta.Response, error) {
	return m.changeAuthorizationServerPolicyLifecycle(ctx, authServerID, policyID, "activate")
}

func (m *ApiSupplement) DeactivateAuthorizationServerPolicy(ctx context.Context, authServerID, policyID string) (*okta.Response, error) {
	return m.changeAuthorizationServerPolicyLifecycle(ctx, authServerID, policyID, "deactivate")
}

func (m *ApiSupplement) changeAuthorizationServerPolicyLifecycle(ctx context.Context, authServerID, policyID, action string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/lifecycle/%s", authServerID, policyID, action)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
