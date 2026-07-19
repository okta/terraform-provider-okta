// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogOutcome represents the LogOutcome schema
type LogOutcome struct {
	// Result of the action
	Result string `json:"result,omitempty"`
	// Reason for the result, for example, `INVALID_CREDENTIALS`
	Reason string `json:"reason,omitempty"`
}
