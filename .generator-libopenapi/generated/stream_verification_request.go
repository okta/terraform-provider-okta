// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// StreamVerificationRequest represents the StreamVerificationRequest schema
type StreamVerificationRequest struct {
	// An arbitrary string that Okta as a transmitter must echo back to the event receiver in the verification event's payload
	State string `json:"state,omitempty"`
	// The ID of the SSF stream configuration
	StreamId string `json:"stream_id"`
}
