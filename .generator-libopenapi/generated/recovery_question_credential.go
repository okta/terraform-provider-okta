// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RecoveryQuestionCredential represents the RecoveryQuestionCredential schema
// Specifies a secret question and answer that's validated (case insensitive) when a user forgets their password or unlocks their account. The answer property is write-only.
type RecoveryQuestionCredential struct {
	// The answer to the recovery question
	Answer string `json:"answer,omitempty"`
	// The recovery question
	Question string `json:"question,omitempty"`
}
