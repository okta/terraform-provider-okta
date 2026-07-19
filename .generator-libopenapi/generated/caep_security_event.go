// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CaepSecurityEvent represents the CaepSecurityEvent schema
type CaepSecurityEvent struct {
	Subject interface{} `json:"subject"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp"`
	// The entity that initiated the event
	InitiatingEntity string `json:"initiating_entity,omitempty"`
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
	ReasonUser interface{} `json:"reason_user,omitempty"`
}
