// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookChannelConfigCreate represents the InlineHookChannelConfigCreate schema
// Properties of the communications channel that are used to contact your external service
type InlineHookChannelConfigCreate struct {
	// An optional list of key/value pairs for headers that you can send with the request to the external service.
	Headers []interface{} `json:"headers,omitempty"`
	// The method of the Okta inline hook request
	Method string `json:"method,omitempty"`
	// The external service endpoint that executes the inline hook handler. It must begin with `https://` and be reachable by Okta. No white space is allowed in the URI.
	Uri string `json:"uri,omitempty"`
}
