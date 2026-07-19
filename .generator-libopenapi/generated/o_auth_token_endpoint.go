// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuthTokenEndpoint represents the OAuthTokenEndpoint schema
// Endpoint for an [OAuth 2.0 Authorization Server (AS)](https://tools.ietf.org/html/rfc6749#page-18)
type OAuthTokenEndpoint struct {
	Binding interface{} `json:"binding,omitempty"`
	// URL of the IdP Authorization Server (AS) token endpoint
	Url string `json:"url,omitempty"`
}
