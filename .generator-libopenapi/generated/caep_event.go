// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CaepEvent represents the CaepEvent schema
type CaepEvent struct {
	ReasonUser interface{} `json:"reason_user,omitempty"`
	Subject interface{} `json:"subject,omitempty"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp,omitempty"`
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
}
