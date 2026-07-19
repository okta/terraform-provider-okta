// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserImportRequestData represents the UserImportRequestData schema
type UserImportRequestData struct {
	// The app user profile being imported
	AppUser map[string]interface{} `json:"appUser,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
	// Provides information on the Okta user profile currently set to be used for the user who is being imported, based on the matching rules and attribute mappings that were applied.
	User map[string]interface{} `json:"user,omitempty"`
	// The object that specifies the default action Okta is set to take
	Action map[string]interface{} `json:"action,omitempty"`
}
