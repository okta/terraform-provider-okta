// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Sso represents the Sso schema
// Supported SSO protocol configurations. You must configure at least one protocol: `oidc` or `saml`
type Sso struct {
	Oidc interface{} `json:"oidc,omitempty"`
	Saml interface{} `json:"saml,omitempty"`
}
