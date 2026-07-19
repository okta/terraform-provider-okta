// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TelephonyRequest represents the TelephonyRequest schema
type TelephonyRequest struct {
	Data interface{} `json:"data,omitempty"`
	// The type of inline hook. The telephony inline hook type is `com.okta.telephony.provider`.
	EventType string `json:"eventType,omitempty"`
	// The type of inline hook request. For example, `com.okta.user.telephony.pre-enrollment`.
	RequestType string `json:"requestType,omitempty"`
	// The ID and URL of the telephony inline hook
	Source string `json:"source,omitempty"`
}
