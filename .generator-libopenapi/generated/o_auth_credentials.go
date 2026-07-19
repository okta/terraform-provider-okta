// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuthCredentials represents the OAuthCredentials schema
// Client authentication credentials for an [OAuth 2.0 Authorization Server](https://tools.ietf.org/html/rfc6749#section-2.3)
type OAuthCredentials struct {
	Client interface{} `json:"client,omitempty"`
	Signing interface{} `json:"signing,omitempty"`
}
