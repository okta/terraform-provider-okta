// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// StreamConfiguration represents the StreamConfiguration schema
type StreamConfiguration struct {
	// The issuer used in security event tokens (SETs). This value is set as `iss` in the claim.  A read-only parameter that is set by the transmitter. If this parameter is included in the request, the va...
	Iss string `json:"iss,omitempty"`
	// The minimum amount of time, in seconds, between two verification requests.  A read-only parameter that is set by the transmitter. If this parameter is included in the request, the value must match ...
	MinVerificationInterval int `json:"min_verification_interval,omitempty"`
	// The audience used in the SET. This value is set as `aud` in the claim.  A read-only parameter that is set by the transmitter. If this parameter is included in the request, the value must match the ...
	Aud interface{} `json:"aud,omitempty"`
	// The ID of the SSF stream configuration
	StreamId string `json:"stream_id,omitempty"`
	Delivery interface{} `json:"delivery"`
	// The events (mapped by the array of event type URIs) that the transmitter actually delivers to the SSF stream.  A read-only parameter that is set by the transmitter. If this parameter is included in...
	EventsDelivered []string `json:"events_delivered,omitempty"`
	// The events (mapped by the array of event type URIs) that the receiver wants to receive
	EventsRequested []string `json:"events_requested"`
	// An array of event type URIs that the transmitter supports.  A read-only parameter that is set by the transmitter. If this parameter is included in the request, the value must match the expected val...
	EventsSupported []string `json:"events_supported,omitempty"`
	// The subject identifier format expected for any SET transmitted.
	Format string `json:"format,omitempty"`
}
