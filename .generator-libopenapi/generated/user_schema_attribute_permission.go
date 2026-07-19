// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserSchemaAttributePermission represents the UserSchemaAttributePermission schema
type UserSchemaAttributePermission struct {
	// Determines whether the principal can view or modify the property
	Action string `json:"action,omitempty"`
	// Security principal
	Principal string `json:"principal,omitempty"`
}
