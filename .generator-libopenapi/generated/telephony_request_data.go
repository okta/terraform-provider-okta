// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TelephonyRequestData represents the TelephonyRequestData schema
type TelephonyRequestData struct {
	// Message profile specifies information about the telephony (sms/voice) message to be sent to the Okta user
	MessageProfile map[string]interface{} `json:"messageProfile,omitempty"`
	// User profile specifies information about the Okta user
	UserProfile map[string]interface{} `json:"userProfile,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}
