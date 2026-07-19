// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// StreamStatus represents the StreamStatus schema
// Status corresponding to the `stream_id` of the SSF stream
type StreamStatus struct {
	// The status of the SSF stream configuration
	Status string `json:"status,omitempty"`
	// The ID of the SSF stream configuration. This corresponds to the value in the query parameter of the request.
	StreamId string `json:"stream_id,omitempty"`
}
