// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BouncesRemoveListError represents the BouncesRemoveListError schema
type BouncesRemoveListError struct {
	// An email address with a validation error
	EmailAddress string `json:"emailAddress,omitempty"`
	// Validation error reason
	Reason string `json:"reason,omitempty"`
}
