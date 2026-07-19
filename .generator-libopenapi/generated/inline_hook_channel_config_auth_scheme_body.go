// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookChannelConfigAuthSchemeBody represents the InlineHookChannelConfigAuthSchemeBody schema
// The authentication scheme to use for this request
type InlineHookChannelConfigAuthSchemeBody struct {
	// The header name for the authorization server
	Key string `json:"key,omitempty"`
	// The authentication scheme type. Supported type&mdash;`HEADER`.
	Type string `json:"type,omitempty"`
	// The header value. This secret value is passed to your external service endpoint. Your external service can check it as a security measure.
	Value string `json:"value,omitempty"`
}
