package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

// FIXME okta-sdk-golang needs to be polymorphic friendly in this manner of passing in a factor instance.
// GetUserFactor fetches a factor for the specified user
func (m *APISupplement) GetUserFactor(ctx context.Context, userId, factorId string, factorInstance okta.Factor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/factors/%v", userId, factorId)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, factorInstance)
}

// FIXME calling undocumented public API
func (m *APISupplement) UpdateUserFactor(ctx context.Context, userId, factorId string, factorInstance okta.Factor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/factors/%v", userId, factorId)
	req, err := m.RequestExecutor.
		WithAccept("application/json").
		WithContentType("application/json").
		NewRequest(http.MethodPut, url, factorInstance)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, factorInstance)
}
