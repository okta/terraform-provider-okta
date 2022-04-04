package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type Factor interface {
	IsUserFactorInstance() bool
}

type SecurityQuestionUserFactor struct {
	Embedded    interface{}                             `json:"_embedded,omitempty"`
	Links       interface{}                             `json:"_links,omitempty"`
	Created     *time.Time                              `json:"created,omitempty"`
	FactorType  string                                  `json:"factorType,omitempty"`
	Id          string                                  `json:"id,omitempty"`
	LastUpdated *time.Time                              `json:"lastUpdated,omitempty"`
	Provider    string                                  `json:"provider,omitempty"`
	Status      string                                  `json:"status,omitempty"`
	Verify      *okta.VerifyFactorRequest               `json:"verify,omitempty"`
	Profile     *okta.SecurityQuestionUserFactorProfile `json:"profile,omitempty"`
}

func (a *SecurityQuestionUserFactor) IsUserFactorInstance() bool {
	return true
}

// EnrollUserFactor enrolls a user with a supported factor.
func (m *APISupplement) EnrollUserFactor(ctx context.Context, userId string, factorInstance Factor, qp *query.Params) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/factors", userId)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.
		WithAccept("application/json").
		WithContentType("application/json").
		NewRequest(http.MethodPost, url, factorInstance)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, factorInstance)
}

// GetUserFactor fetches a factor for the specified user
func (m *APISupplement) GetUserFactor(ctx context.Context, userId, factorId string, factorInstance Factor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/factors/%v", userId, factorId)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, factorInstance)
}

func (m *APISupplement) UpdateUserFactor(ctx context.Context, userId, factorId string, factorInstance Factor) (*okta.Response, error) {
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
