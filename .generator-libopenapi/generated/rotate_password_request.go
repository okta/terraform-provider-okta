// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RotatePasswordRequest represents the RotatePasswordRequest schema
// Rotate password request for the privileged resource
type RotatePasswordRequest struct {
	// The password associated with the privileged resource
	Password string `json:"password"`
	// The version ID of the password secret from the OPA vault
	SecretVersionId string `json:"secretVersionId"`
}
