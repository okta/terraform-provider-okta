package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type Subscription struct {
	NotificationType string      `json:"notificationType"`
	Channels         []string    `json:"channels"`
	Status           string      `json:"status"`
	Links            interface{} `json:"links"`
}

func (m *APISupplement) GetRoleTypeSubscription(ctx context.Context, roleType, notificationType string) (*Subscription, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/roles/%s/subscriptions/%s", roleType, notificationType)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	var subscription *Subscription
	resp, err := re.Do(ctx, req, &subscription)
	if err != nil {
		return nil, resp, err
	}
	return subscription, resp, nil
}

func (m *APISupplement) RoleTypeSubscribe(ctx context.Context, roleType, notificationType string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/roles/%s/subscriptions/%s/subscribe", roleType, notificationType)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}

func (m *APISupplement) RoleTypeUnsubscribe(ctx context.Context, roleType, notificationType string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/roles/%s/subscriptions/%s/unsubscribe", roleType, notificationType)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
