// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEvent represents the SecurityEvent schema
type SecurityEvent struct {
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp"`
	Subject interface{} `json:"subject"`
}
