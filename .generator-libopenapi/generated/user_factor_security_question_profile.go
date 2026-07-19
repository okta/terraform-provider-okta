// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorSecurityQuestionProfile represents the UserFactorSecurityQuestionProfile schema
type UserFactorSecurityQuestionProfile struct {
	// Answer to the question
	Answer string `json:"answer,omitempty"`
	// Unique key for the question
	Question string `json:"question,omitempty"`
	// Human-readable text that's displayed to the user
	QuestionText string `json:"questionText,omitempty"`
}
