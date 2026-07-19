// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CaepDeviceComplianceChangeEvent represents the CaepDeviceComplianceChangeEvent schema
// The subject's device compliance was revoked
type CaepDeviceComplianceChangeEvent struct {
	Subject interface{} `json:"subject"`
	// Current device compliance status
	CurrentStatus string `json:"current_status"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp"`
	// The entity that initiated the event
	InitiatingEntity string `json:"initiating_entity,omitempty"`
	// Previous device compliance status
	PreviousStatus string `json:"previous_status"`
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
	ReasonUser interface{} `json:"reason_user,omitempty"`
}
