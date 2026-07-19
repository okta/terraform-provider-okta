// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// tempPassword represents the tempPassword schema
type tempPassword struct {
	// A temporary password that the user can sign in with. This is only returned when expiring a password with a temporary password.
	TempPassword string `json:"tempPassword,omitempty"`
}
