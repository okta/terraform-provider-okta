// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookOAuthPrivateKeyJwtConfig represents the InlineHookOAuthPrivateKeyJwtConfig schema
type InlineHookOAuthPrivateKeyJwtConfig struct {
	// Not applicable. Must be `null`.
	AuthScheme string `json:"authScheme,omitempty"`
	// An ID value of the hook key pair generated from the [Hook Keys API](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/HookKey/#tag/HookKey)
	HookKeyId string `json:"hookKeyId,omitempty"`
	// The method of the Okta inline hook request. Only accepts `POST`.
	Method string `json:"method,omitempty"`
}
