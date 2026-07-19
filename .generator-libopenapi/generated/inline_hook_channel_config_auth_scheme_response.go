// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookChannelConfigAuthSchemeResponse represents the InlineHookChannelConfigAuthSchemeResponse schema
// The authentication scheme to use for this request
type InlineHookChannelConfigAuthSchemeResponse struct {
	// The header name for the authorization server
	Key string `json:"key,omitempty"`
	// The authentication scheme type. Supported type&mdash;`HEADER`
	Type string `json:"type,omitempty"`
}
