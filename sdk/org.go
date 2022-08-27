package sdk

import (
	"context"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type OktaOrganization struct {
	Id       string      `json:"id"`
	Pipeline string      `json:"pipeline"`
	Links    interface{} `json:"_links,omitempty"`
	Settings interface{} `json:"settings,omitempty"`
}

// GetWellKnownOktaOrganization calls GET /.well-known/okta-organization that
// surfaces information about the org including if it is OIE or Classic
// (pipeline=v1 is Classic, pipeline=idx is OIE)
//
// NOTE: Devs at Okta are negotiating internally to recognize the endpoint as
// public and will be in a coming release of okta-sdk-golang and documented at
// developer.okta.com .
//
// TODO: remove this custom code with well known okta organization is in okta-sdk-golang
func (m *APISupplement) GetWellKnownOktaOrganization(ctx context.Context) (*OktaOrganization, *okta.Response, error) {
	url := "/.well-known/okta-organization"
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var oktaOrganization *OktaOrganization
	resp, err := m.RequestExecutor.Do(ctx, req, &oktaOrganization)
	if err != nil {
		return nil, resp, err
	}
	return oktaOrganization, resp, nil
}
