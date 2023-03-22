package sdk

import (
	"context"
	"fmt"
	"time"
)

type SessionResource resource

type Session struct {
	Links                    interface{}                    `json:"_links,omitempty"`
	Amr                      []*SessionAuthenticationMethod `json:"amr,omitempty"`
	CreatedAt                *time.Time                     `json:"createdAt,omitempty"`
	ExpiresAt                *time.Time                     `json:"expiresAt,omitempty"`
	Id                       string                         `json:"id,omitempty"`
	Idp                      *SessionIdentityProvider       `json:"idp,omitempty"`
	LastFactorVerification   *time.Time                     `json:"lastFactorVerification,omitempty"`
	LastPasswordVerification *time.Time                     `json:"lastPasswordVerification,omitempty"`
	Login                    string                         `json:"login,omitempty"`
	Status                   string                         `json:"status,omitempty"`
	UserId                   string                         `json:"userId,omitempty"`
}

// Get details about a session.
func (m *SessionResource) GetSession(ctx context.Context, sessionId string) (*Session, *Response, error) {
	url := fmt.Sprintf("/api/v1/sessions/%v", sessionId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var session *Session

	resp, err := rq.Do(ctx, req, &session)
	if err != nil {
		return nil, resp, err
	}

	return session, resp, nil
}

func (m *SessionResource) EndSession(ctx context.Context, sessionId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/sessions/%v", sessionId)

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

// Creates a new session for a user with a valid session token. Use this API if, for example, you want to set the session cookie yourself instead of allowing Okta to set it, or want to hold the session ID in order to delete a session via the API instead of visiting the logout URL.
func (m *SessionResource) CreateSession(ctx context.Context, body CreateSessionRequest) (*Session, *Response, error) {
	url := fmt.Sprintf("/api/v1/sessions")

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var session *Session

	resp, err := rq.Do(ctx, req, &session)
	if err != nil {
		return nil, resp, err
	}

	return session, resp, nil
}

func (m *SessionResource) RefreshSession(ctx context.Context, sessionId string) (*Session, *Response, error) {
	url := fmt.Sprintf("/api/v1/sessions/%v/lifecycle/refresh", sessionId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var session *Session

	resp, err := rq.Do(ctx, req, &session)
	if err != nil {
		return nil, resp, err
	}

	return session, resp, nil
}
