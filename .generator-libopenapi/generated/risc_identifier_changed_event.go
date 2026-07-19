// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RiscIdentifierChangedEvent represents the RiscIdentifierChangedEvent schema
// The subject's identifier has changed, which is either an email address or a phone number change
type RiscIdentifierChangedEvent struct {
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp"`
	// The new identifier value
	New-value string `json:"new-value,omitempty"`
	Subject interface{} `json:"subject"`
}
