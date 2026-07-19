// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SocialAuthToken represents the SocialAuthToken schema
// The social authentication token object provides the tokens and associated metadata provided by social providers during social authentication.
type SocialAuthToken struct {
	// The type of token defined by the [OAuth Token Exchange Spec](https://tools.ietf.org/html/draft-ietf-oauth-token-exchange-07#section-3)
	TokenType string `json:"tokenType,omitempty"`
	ExpiresAt interface{} `json:"expiresAt,omitempty"`
	// Unique identifier for the token
	ID string `json:"id,omitempty"`
	// The scopes that the token is good for
	Scopes []string `json:"scopes,omitempty"`
	// The raw token
	Token string `json:"token,omitempty"`
	// The token authentication scheme as defined by the social provider
	TokenAuthScheme string `json:"tokenAuthScheme,omitempty"`
}
