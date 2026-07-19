// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyRecoveryFactors represents the PasswordPolicyRecoveryFactors schema
// Settings for the factors that can be used for recovery
type PasswordPolicyRecoveryFactors struct {
	// Okta voice call
	OktaCall interface{} `json:"okta_call,omitempty"`
	// Okta email
	OktaEmail interface{} `json:"okta_email,omitempty"`
	// Okta SMS
	OktaSms interface{} `json:"okta_sms,omitempty"`
	// Okta security question
	RecoveryQuestion interface{} `json:"recovery_question,omitempty"`
}
