package okta

import (
	"fmt"
	"time"
)

type TrustedOriginsService service

func (p *TrustedOriginsService) TrustedOrigin() TrustedOrigin {
	return TrustedOrigin{}
}

type TrustedOrigin struct {
	ID            string              `json:"id,omitempty"`
	Status        string              `json:"status,omitempty"`
	Name          string              `json:"name,omitempty"`
	Origin        string              `json:"origin,omitempty"`
	Scopes        []map[string]string `json:"scopes,omitempty"`
	Created       *time.Time          `json:"created,omitempty"`
	CreatedBy     string              `json:"createdBy,omitempty"`
	LastUpdated   *time.Time          `json:"lastUpdated,omitempty"`
	LastUpdatedBy string              `json:"lastUpdated,omitempty"`
	Links         *TrustedOriginLinks `json:"_links,omitempty"`
}

type TrustedOriginDeactive struct {
	Href  string              `json:"href,omitempty"`
	Hints *TrustedOriginHints `json:"hints,omitempty"`
}

type TrustedOriginHints struct {
	Allow []string `json:"allow,omitempty"`
}

type TrustedOriginLinks struct {
	Self       *TrustedOriginSelf     `json:"self,omitempty"`
	Deactivate *TrustedOriginDeactive `json:"deactive,omitempty"`
}

type TrustedOriginSelf struct {
	Href  string              `json:"href,omitempty"`
	Hints *TrustedOriginHints `json:"hints,omitempty"`
}

// GetTrustedOrigin: Get a Trusted Origin entry
// Requires TrustedOrigins ID from TrustedOrigins object
func (p *TrustedOriginsService) GetTrustedOrigin(id string) (*TrustedOrigin, *Response, error) {
	u := fmt.Sprintf("trustedOrigins/%v", id)
	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	trustedOrigin := new(TrustedOrigin)
	resp, err := p.client.Do(req, trustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return trustedOrigin, resp, err
}

// CreateTrustedOrigin: Create a Trusted Origin
// You must pass in the Trusted Origin object created from the desired input trustedOrigin
func (p *TrustedOriginsService) CreateTrustedOrigin(trustedOrigin interface{}) (*TrustedOrigin, *Response, error) {
	u := fmt.Sprintf("trustedOrigins")
	req, err := p.client.NewRequest("POST", u, trustedOrigin)

	if err != nil {
		return nil, nil, err
	}

	newTrustedOrigin := new(TrustedOrigin)

	resp, err := p.client.Do(req, newTrustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return newTrustedOrigin, resp, err
}

// UpdateTrustedOrigin: Update a Trusted Origin
// Requires TrustedOrigin ID from TrustedOrigin object & TrustedOrigin object from the desired input policy
func (p *TrustedOriginsService) UpdateTrustedOrigin(id string, trustedOrigin interface{}) (*TrustedOrigin, *Response, error) {
	u := fmt.Sprintf("trustedOrigins/%v", id)
	req, err := p.client.NewRequest("PUT", u, trustedOrigin)
	if err != nil {
		return nil, nil, err
	}

	updateTrustedOrigin := new(TrustedOrigin)
	resp, err := p.client.Do(req, updateTrustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return updateTrustedOrigin, resp, err
}

// DeleteTrustedOrigin: Delete a Trusted Origin
// Requires TrustedOrigin ID from TrustedOrigin object
func (p *TrustedOriginsService) DeleteTrustedOrigin(id string) (*Response, error) {
	u := fmt.Sprintf("trustedOrigins/%v", id)
	req, err := p.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// ActivateTrustedOrigin: Activate/Deactivate a Trusted Origin
// Requires TrustedOrigin ID from TrustedOrigin object and a boolean to activate or deactivate
func (p *TrustedOriginsService) ActivateTrustedOrigin(id string, activate bool) (*Response, error) {
	var a string

	if activate {
		a = "activate"
	} else {
		a = "deactivate"
	}

	u := fmt.Sprintf("trustedOrigins/%v/lifecycle/%v", id, a)
	req, err := p.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// ListTrustedOrigins: Lists all Trusted Origins from an Okta Account
func (p *TrustedOriginsService) ListTrustedOrigins() (*Response, error) {
	u := fmt.Sprintf("trustedOrigins")
	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
