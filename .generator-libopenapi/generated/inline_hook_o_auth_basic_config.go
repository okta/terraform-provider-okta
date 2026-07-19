// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookOAuthBasicConfig represents the InlineHookOAuthBasicConfig schema
type InlineHookOAuthBasicConfig struct {
	// Include the scopes that allow you to perform the actions on the hook endpoint that you want to access
	Scope string `json:"scope,omitempty"`
	// The URI where inline hooks can exchange an authorization code for access and refresh tokens
	TokenUrl string `json:"tokenUrl,omitempty"`
	AuthType string `json:"authType,omitempty"`
	// A publicly exposed string provided by the service that's used to identify the OAuth app and build authorization URLs
	ClientId string `json:"clientId,omitempty"`
}
