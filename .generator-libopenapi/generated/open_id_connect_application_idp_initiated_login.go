// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OpenIdConnectApplicationIdpInitiatedLogin represents the OpenIdConnectApplicationIdpInitiatedLogin schema
// The type of IdP-initiated sign-in flow that the client supports
type OpenIdConnectApplicationIdpInitiatedLogin struct {
	// The scopes to use for the request when `mode` is `OKTA`
	DefaultScope []string `json:"default_scope,omitempty"`
	// The mode to use for the IdP-initiated sign-in flow. For `OKTA` or `SPEC` modes, the client must have an `initiate_login_uri` registered. > **Note:** For web and SPA apps, if the mode is `SPEC` or `...
	Mode string `json:"mode"`
}
