// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EventHookChannelConfig represents the EventHookChannelConfig schema
type EventHookChannelConfig struct {
	// Optional list of key/value pairs for headers that can be sent with the request to the external service. For example, `X-Other-Header` is an example of an optional header, with a value of `my-header...
	Headers []interface{} `json:"headers,omitempty"`
	// The method of the Okta event hook request
	Method string `json:"method,omitempty"`
	// The external service endpoint called to execute the event hook handler
	Uri string `json:"uri"`
	AuthScheme interface{} `json:"authScheme,omitempty"`
}
