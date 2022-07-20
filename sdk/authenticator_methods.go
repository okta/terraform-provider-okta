package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type authenticator_method struct {
	Links  interface{} `json:"_links,omitempty"`
	Type   string      `json:"type,omitempty"`
	Status string      `json:"status,omitempty"`
}

// GetAppSignOnPolicyRule gets a policy rule.
func (m *APISupplement) GetAuthenticatorMethodStatus(ctx context.Context, authenticatorID, method string) (string, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v/methods/%v", authenticatorID, method)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	var authenticatorMethod authenticator_method
	_, err = m.RequestExecutor.Do(ctx, req, &authenticatorMethod)
	if err != nil {
		return "", err
	}
	return authenticatorMethod.Status, nil
}

// ActivateAuthenticatorMethod deactivates the Authenticators supplied method.
func (m *APISupplement) ActivateAuthenticatorMethod(ctx context.Context, authenticatorID, method string) error {
	status, err := m.GetAuthenticatorMethodStatus(ctx, authenticatorID, method)
	if err != nil {
		return err
	}
	if status != "ACTIVE" {
		return m.lifecycleChangeAuthenticatorMethod(ctx, authenticatorID, method, "activate", 1)
	}
	return nil
}

// DeactivateAuthenticatorMethod deactivates the Authenticators supplied method.
func (m *APISupplement) DeactivateAuthenticatorMethod(ctx context.Context, authenticatorID, method string) error {
	status, err := m.GetAuthenticatorMethodStatus(ctx, authenticatorID, method)
	if err != nil {
		return err
	}
	if status != "INACTIVE" {
		return m.lifecycleChangeAuthenticatorMethod(ctx, authenticatorID, method, "deactivate", 1)
	}
	return nil
}

func (m *APISupplement) lifecycleChangeAuthenticatorMethod(ctx context.Context, authenticatorID, method, action string, attempt int) error {
	url := fmt.Sprintf("/api/v1/authenticators/%s/methods/%s/lifecycle/%s", authenticatorID, method, action)
	req, err := m.RequestExecutor.
		WithAccept("application/json").
		WithContentType("application/json").
		NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	_, err = m.RequestExecutor.Do(ctx, req, nil)
	if err != nil {
		// while testing received a few 501s for /api/v1/authenticators/%s/methods/%s/lifecycle/%s - could be possible lag for this endpoint
		// to be available/consistent on backend. Retry if so
		if err.Error() == "the API returned an error: Unsupported operation." {
			if attempt > 5 {
				return err
			}
			time.Sleep(time.Duration(attempt * 1000000000))
			return m.lifecycleChangeAuthenticatorMethod(ctx, authenticatorID, method, action, attempt+1)
		}
		return err
	}
	return err
}
