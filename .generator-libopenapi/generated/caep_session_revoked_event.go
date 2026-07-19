// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CaepSessionRevokedEvent represents the CaepSessionRevokedEvent schema
// The session of the subject was revoked
type CaepSessionRevokedEvent struct {
	ReasonUser interface{} `json:"reason_user,omitempty"`
	// Current IP of the session
	CurrentIp string `json:"current_ip,omitempty"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp"`
	// The entity that initiated the event
	InitiatingEntity string `json:"initiating_entity,omitempty"`
	// Last known user agent of the session
	LastKnownUserAgent string `json:"last_known_user_agent,omitempty"`
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
	Subject interface{} `json:"subject"`
	// Current user agent of the session
	CurrentUserAgent string `json:"current_user_agent,omitempty"`
	// Last known IP of the session
	LastKnownIp string `json:"last_known_ip,omitempty"`
}
