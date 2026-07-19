// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordImportRequest represents the PasswordImportRequest schema
type PasswordImportRequest struct {
	Data interface{} `json:"data,omitempty"`
	// The type of inline hook. The password import inline hook type is `com.okta.user.credential.password.import`.
	EventType string `json:"eventType,omitempty"`
	// The ID and URL of the password import inline hook
	Source string `json:"source,omitempty"`
}
