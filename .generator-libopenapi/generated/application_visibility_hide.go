// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationVisibilityHide represents the ApplicationVisibilityHide schema
// Hides the app for specific end-user apps
type ApplicationVisibilityHide struct {
	// Okta Mobile for iOS or Android (pre-dates Android)
	IOS bool `json:"iOS,omitempty"`
	// Okta End-User Dashboard on a web browser
	Web bool `json:"web,omitempty"`
}
