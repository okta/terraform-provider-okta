// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventTokenError represents the SecurityEventTokenError schema
// Error object thrown when parsing the security event token
type SecurityEventTokenError struct {
	// Describes the error > **Note:** SET claim fields with underscores (snake case) are presented in camelcase. For example, `previous_status` appears as `previousStatus`.
	Description string `json:"description,omitempty"`
	// A code that describes the category of the error
	Err string `json:"err,omitempty"`
}
