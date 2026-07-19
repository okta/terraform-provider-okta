// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OpenIdConnectApplicationSettingsRefreshToken represents the OpenIdConnectApplicationSettingsRefreshToken schema
// Refresh token configuration for an OAuth 2.0 client  When you create or update an OAuth 2.0 client, you can configure refresh token rotation by setting the `rotation_type` and `leeway` properties. ...
type OpenIdConnectApplicationSettingsRefreshToken struct {
	// The leeway, in seconds, allowed for the OAuth 2.0 client. After the refresh token is rotated, the previous token remains valid for the specified period of time so clients can get the new token.  > ...
	Leeway int `json:"leeway,omitempty"`
	RotationType interface{} `json:"rotation_type"`
}
