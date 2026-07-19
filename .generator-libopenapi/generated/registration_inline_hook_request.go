// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RegistrationInlineHookRequest represents the RegistrationInlineHookRequest schema
// Registration inline hook request
type RegistrationInlineHookRequest struct {
	// The type of inline hook. The registration inline hook type is `com.okta.user.pre-registration`.
	EventType string `json:"eventType,omitempty"`
	RequestType interface{} `json:"requestType,omitempty"`
	// The ID of the registration inline hook
	Source string `json:"source,omitempty"`
}
