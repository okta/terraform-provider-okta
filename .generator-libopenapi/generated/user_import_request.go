// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserImportRequest represents the UserImportRequest schema
type UserImportRequest struct {
	Data interface{} `json:"data,omitempty"`
	// The type of inline hook. The user import inline hook type is `com.okta.import.transform`.
	EventType string `json:"eventType,omitempty"`
	// The ID of the user import inline hook
	Source string `json:"source,omitempty"`
}
