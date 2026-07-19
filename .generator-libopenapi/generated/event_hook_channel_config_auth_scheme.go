// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EventHookChannelConfigAuthScheme represents the EventHookChannelConfigAuthScheme schema
// The authentication scheme used for this request.  To use Basic Auth for authentication, set `type` to `HEADER`, `key` to `Authorization`, and `value` to the Base64-encoded string of "username:passw...
type EventHookChannelConfigAuthScheme struct {
	// The name for the authorization header
	Key string `json:"key,omitempty"`
	Type interface{} `json:"type,omitempty"`
	// The header value. This secret key is passed to your external service endpoint for security verification. This property is not returned in the response.
	Value string `json:"value,omitempty"`
}
