package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type AppUserResource resource

type AppUser struct {
	Embedded        interface{}         `json:"_embedded,omitempty"`
	Links           interface{}         `json:"_links,omitempty"`
	Created         *time.Time          `json:"created,omitempty"`
	Credentials     *AppUserCredentials `json:"credentials,omitempty"`
	ExternalId      string              `json:"externalId,omitempty"`
	Id              string              `json:"id,omitempty"`
	LastSync        *time.Time          `json:"lastSync,omitempty"`
	LastUpdated     *time.Time          `json:"lastUpdated,omitempty"`
	PasswordChanged *time.Time          `json:"passwordChanged,omitempty"`
	Profile         interface{}         `json:"profile,omitempty"`
	Scope           string              `json:"scope,omitempty"`
	Status          string              `json:"status,omitempty"`
	StatusChanged   *time.Time          `json:"statusChanged,omitempty"`
	SyncState       string              `json:"syncState,omitempty"`
}

// Updates a user&#x27;s profile for an application
func (m *AppUserResource) UpdateApplicationUser(ctx context.Context, appId string, userId string, body AppUser) (*AppUser, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users/%v", appId, userId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var appUser *AppUser

	resp, err := rq.Do(ctx, req, &appUser)
	if err != nil {
		return nil, resp, err
	}

	return appUser, resp, nil
}

// Removes an assignment for a user from an application.
func (m *AppUserResource) DeleteApplicationUser(ctx context.Context, appId string, userId string, qp *query.Params) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users/%v", appId, userId)
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
