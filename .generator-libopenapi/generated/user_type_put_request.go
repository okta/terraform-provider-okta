// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserTypePutRequest represents the UserTypePutRequest schema
type UserTypePutRequest struct {
	// The human-readable description of the user type
	Description string `json:"description"`
	// The human-readable name of the user type
	DisplayName string `json:"displayName"`
	// The name of the existing type
	Name string `json:"name"`
}
