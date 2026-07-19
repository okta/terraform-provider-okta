// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaUserRiskChangeEvent represents the OktaUserRiskChangeEvent schema
// The user risk level changed
type OktaUserRiskChangeEvent struct {
	Subject interface{} `json:"subject"`
	// Current risk level of the user
	CurrentLevel string `json:"current_level"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp"`
	// The entity that initiated the event
	InitiatingEntity string `json:"initiating_entity,omitempty"`
	// Previous risk level of the user
	PreviousLevel string `json:"previous_level"`
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
	ReasonUser interface{} `json:"reason_user,omitempty"`
}
