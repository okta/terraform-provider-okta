// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SsfTransmitterCaepSessionRevokedEvent represents the SsfTransmitterCaepSessionRevokedEvent schema
// The session of the subject was revoked
type SsfTransmitterCaepSessionRevokedEvent struct {
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
	ReasonUser interface{} `json:"reason_user,omitempty"`
	Subject interface{} `json:"subject,omitempty"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp,omitempty"`
}
