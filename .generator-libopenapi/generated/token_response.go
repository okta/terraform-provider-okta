// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TokenResponse represents the TokenResponse schema
type TokenResponse struct {
	// The expiration time of the access token in seconds.
	ExpiresIn int `json:"expires_in,omitempty"`
	// An ID token. This is returned if the `openid` scope is granted.
	IdToken string `json:"id_token,omitempty"`
	IssuedTokenType interface{} `json:"issued_token_type,omitempty"`
	// An opaque refresh token. This is returned if the `offline_access` scope is granted.
	RefreshToken string `json:"refresh_token,omitempty"`
	// The scopes contained in the access token.
	Scope string `json:"scope,omitempty"`
	TokenType interface{} `json:"token_type,omitempty"`
	// An access token.
	AccessToken string `json:"access_token,omitempty"`
	// An opaque device secret. This is returned if the `device_sso` scope is granted.
	DeviceSecret string `json:"device_secret,omitempty"`
}
