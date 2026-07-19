// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// StreamConfigurationCreateRequest represents the StreamConfigurationCreateRequest schema
type StreamConfigurationCreateRequest struct {
	// The events (mapped by the array of event type URIs) that the receiver wants to receive
	EventsRequested []string `json:"events_requested"`
	// The subject identifier format expected for any SET transmitted.
	Format string `json:"format,omitempty"`
	Delivery interface{} `json:"delivery"`
}
