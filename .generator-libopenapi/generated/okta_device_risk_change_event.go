// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaDeviceRiskChangeEvent represents the OktaDeviceRiskChangeEvent schema
// The device risk level changed
type OktaDeviceRiskChangeEvent struct {
	// Current risk level of the device
	CurrentLevel string `json:"current_level"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp"`
	// The entity that initiated the event
	InitiatingEntity string `json:"initiating_entity,omitempty"`
	// Previous risk level of the device
	PreviousLevel string `json:"previous_level"`
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
	ReasonUser interface{} `json:"reason_user,omitempty"`
	Subject interface{} `json:"subject"`
}
