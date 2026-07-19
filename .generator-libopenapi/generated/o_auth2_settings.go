// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2Settings represents the OAuth2Settings schema
// OAuth 2.0 configuration used for authType `OAUTH2`
type OAuth2Settings struct {
	// The URL to the authorization server's authorization endpoint
	AuthorizeEndpoint string `json:"authorizeEndpoint"`
	// The OAuth 2.0 client identifier
	ClientId string `json:"clientId"`
	// The OAuth 2.0 client secret
	ClientSecret string `json:"clientSecret"`
	// The public key in JWK format. Returned when the OAuth authentication method is `PRIVATE_KEY_JWT`.
	PublicKey map[string]interface{} `json:"publicKey,omitempty"`
	// List of OAuth 2.0 scopes
	Scopes []string `json:"scopes,omitempty"`
	// The URL to the authorization server's token endpoint
	TokenEndpoint string `json:"tokenEndpoint"`
}
