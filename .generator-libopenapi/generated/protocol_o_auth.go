// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProtocolOAuth represents the ProtocolOAuth schema
// Protocol settings for authentication using the [OAuth 2.0 Authorization Code flow](https://tools.ietf.org/html/rfc6749#section-4.1)
type ProtocolOAuth struct {
	Credentials interface{} `json:"credentials,omitempty"`
	Endpoints interface{} `json:"endpoints,omitempty"`
	Scopes interface{} `json:"scopes,omitempty"`
	// OAuth 2.0 Authorization Code flow
	Type string `json:"type,omitempty"`
}
