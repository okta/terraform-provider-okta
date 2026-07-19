// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaIpChangeEvent represents the OktaIpChangeEvent schema
// IP changed for the subject's session
type OktaIpChangeEvent struct {
	ReasonUser interface{} `json:"reason_user,omitempty"`
	Subject interface{} `json:"subject"`
	// Current IP address of the subject
	CurrentIpAddress string `json:"current_ip_address"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp"`
	// The entity that initiated the event
	InitiatingEntity string `json:"initiating_entity,omitempty"`
	// Previous IP address of the subject
	PreviousIpAddress string `json:"previous_ip_address"`
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
}
