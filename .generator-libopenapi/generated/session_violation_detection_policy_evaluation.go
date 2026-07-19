// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SessionViolationDetectionPolicyEvaluation represents the SessionViolationDetectionPolicyEvaluation schema
// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Used to control evaluation of the session sign-on policies
type SessionViolationDetectionPolicyEvaluation struct {
	// When true, the sign-on policies of the session are evaluated
	Enabled bool `json:"enabled,omitempty"`
}
