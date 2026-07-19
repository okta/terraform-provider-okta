// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlSloEndpoint represents the SamlSloEndpoint schema
// IdP's `SingleLogoutService` endpoint where Okta sends a `<LogoutRequest>` message
type SamlSloEndpoint struct {
	Binding interface{} `json:"binding,omitempty"`
	// URL of the binding-specific IdP endpoint where Okta sends a `<LogoutRequest>`
	Url string `json:"url,omitempty"`
}
