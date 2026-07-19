// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationSettingsNotes represents the ApplicationSettingsNotes schema
// App notes visible to either the admin or end user
type ApplicationSettingsNotes struct {
	// An app message that's visible to admins
	Admin string `json:"admin,omitempty"`
	// A message that's visible in the End-User Dashboard
	Enduser string `json:"enduser,omitempty"`
}
