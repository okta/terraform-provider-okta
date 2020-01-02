package sdk

import (
	"github.com/okta/okta-sdk-golang/okta"
)

type (
	OpenIdConnectApplication struct {
		okta.OpenIdConnectApplication
		Settings *OpenIdConnectApplicationSettings `json:"settings,omitempty"`
	}

	OpenIdConnectApplicationSettings struct {
		okta.OpenIdConnectApplicationSettingsClient
		OauthClient *OpenIdConnectApplicationSettingsClient `json:"oauthClient,omitempty"`
	}

	OpenIdConnectApplicationSettingsClient struct {
		okta.OpenIdConnectApplicationSettingsClient
		JWKS *JWKS `json:"jwks,omitempty"`
	}

	JWKS struct {
		Keys []*JWK `json:"keys,omitempty"`
	}

	JWK struct {
		Type     string `json:"kty,omitempty"`
		ID       string `json:"kid,omitempty"`
		Exponent string `json:"e,omitempty"`
		Modulus  string `json:"n,omitempty"`
	}
)

func NewOpenIdConnectApplication() *OpenIdConnectApplication {
	app := &OpenIdConnectApplication{}
	app.Name = "oidc_client"
	app.SignOnMode = "OPENID_CONNECT"
	return app
}

func (a *OpenIdConnectApplication) IsApplicationInstance() bool {
	return true
}
