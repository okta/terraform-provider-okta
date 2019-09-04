package okta

import (
	"fmt"
	"time"
)

type IdentityProvidersService service

func (p *IdentityProvidersService) IdentityProvider() IdentityProvider {
	return IdentityProvider{}
}

type AccountLink struct {
	Filter string `json:"filter,omitempty"`
	Action string `json:"action,omitempty"`
}

type Authorization struct {
	Url     string `json:"url,omitempty"`
	Binding string `json:"binding,omitempty"`
}

type Authorize struct {
	Href      string `json:"href,omitempty"`
	Templated bool   `json:"templated,omitempty"`
	Hints     *Hints `json:"hints,omitempty"`
}

type IdpClient struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

type ClientRedirectUri struct {
	Href  string `json:"href,omitempty"`
	Hints *Hints `json:"hints,omitempty"`
}

type Conditions struct {
	Deprovisioned *Deprovisioned `json:"deprovisioned,omitempty"`
	Suspended     *Suspended     `json:"suspended,omitempty"`
}

type Credentials struct {
	Client *IdpClient `json:"client,omitempty"`
}

type Deprovisioned struct {
	Action string `json:"action,omitempty"`
}

type Endpoints struct {
	Authorization *Authorization `json:"authorization,omitempty"`
	Token         *Token         `json:"token,omitempty"`
}

type IdpGroups struct {
	Action      string   `json:"action,omitempty"`
	Assignments []string `json:"assignments,omitempty"`
}

type Hints struct {
	Allow []string `json:"allow,omitempty"`
}

// Note - time.Time fields are pointers due to the issue described at link below
// https://stackoverflow.com/questions/32643815/golang-json-omitempty-with-time-time-field
type IdentityProvider struct {
	ID          string     `json:"id,omitempty"`
	Type        string     `json:"type,omitempty"`
	Status      string     `json:"status,omitempty"`
	IssuerMode  string     `json:"issuerMode,omitempty"`
	Name        string     `json:"name,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Protocol    *Protocol  `json:"protocol,omitempty"`
	Policy      *IdpPolicy `json:"policy,omitempty"`
	Links       *IdpLinks  `json:"_links,omitempty"`
}

type IdpLinks struct {
	Authorize         *Authorize         `json:"authorize,omitempty"`
	ClientRedirectUri *ClientRedirectUri `json:"clientRedirectUri,omitempty"`
}

type IdpPolicy struct {
	Provisioning *Provisioning `json:"provisioning,omitempty"`
	AccountLink  *AccountLink  `json:"accountLink,omitempty"`
	Subject      *Subject      `json:"subject,omitempty"`
	MaxClockSkew int           `json:"maxClockSkew,omitempty"`
}

type Protocol struct {
	Type        string       `json:"type,omitempty"`
	Endpoints   *Endpoints   `json:"endpoints,omitempty"`
	Scopes      []string     `json:"scopes,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty"`
}

type Provisioning struct {
	Action        string      `json:"action,omitempty"`
	ProfileMaster bool        `json:"profileMaster,omitempty"`
	Groups        *IdpGroups  `json:"groups,omitempty"`
	Conditions    *Conditions `json:"conditions,omitempty"`
}

type Suspended struct {
	Action string `json:"action,omitempty"`
}

type Subject struct {
	UserNameTemplate *UserNameTemplate `json:"userNameTemplate,omitempty"`
	Filter           string            `json:"filter,omitempty"`
	MatchType        string            `json:"matchType,omitempty"`
}

type Token struct {
	Url     string `json:"url,omitempty"`
	Binding string `json:"binding,omitempty"`
}

type UserNameTemplate struct {
	Template string `json:"template,omitempty"`
}

// GetIdentityProvider: Get an Identity Provider
// Requires IdentityProvider ID from IdentityProvider object
func (p *IdentityProvidersService) GetIdentityProvider(id string) (*IdentityProvider, *Response, error) {
	u := fmt.Sprintf("idps/%v", id)
	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	idp := new(IdentityProvider)
	resp, err := p.client.Do(req, idp)
	if err != nil {
		return nil, resp, err
	}

	return idp, resp, err
}

// CreateIdentityProvider: Create an Identity Provider
// You must pass in the IdentityProvider object created from the desired input IdentityProvider
func (p *IdentityProvidersService) CreateIdentityProvider(idp interface{}) (*IdentityProvider, *Response, error) {
	u := fmt.Sprintf("idps")
	req, err := p.client.NewRequest("POST", u, idp)

	if err != nil {
		return nil, nil, err
	}

	newIdp := new(IdentityProvider)

	resp, err := p.client.Do(req, newIdp)
	if err != nil {
		return nil, resp, err
	}

	return newIdp, resp, err
}

// UpdateIdentityProvider: Update an Identity Provider
// Requires IdentityProvider ID from IdentityProvider object & IdentityProvider object from the desired input IdentityProvider
func (p *IdentityProvidersService) UpdateIdentityProvider(id string, idp interface{}) (*IdentityProvider, *Response, error) {
	u := fmt.Sprintf("idps/%v", id)
	req, err := p.client.NewRequest("PUT", u, idp)
	if err != nil {
		return nil, nil, err
	}

	updateIdentityProvider := new(IdentityProvider)
	resp, err := p.client.Do(req, updateIdentityProvider)
	if err != nil {
		return nil, resp, err
	}

	return updateIdentityProvider, resp, err
}

// DeleteIdentityProvider: Delete an Identity Provider
// Requires IdentityProvider ID from IdentityProvider object
func (p *IdentityProvidersService) DeleteIdentityProvider(id string) (*Response, error) {
	u := fmt.Sprintf("idps/%v", id)
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
