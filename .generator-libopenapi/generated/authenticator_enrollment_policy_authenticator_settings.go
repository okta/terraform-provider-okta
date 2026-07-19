// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticatorEnrollmentPolicyAuthenticatorSettings represents the AuthenticatorEnrollmentPolicyAuthenticatorSettings schema
type AuthenticatorEnrollmentPolicyAuthenticatorSettings struct {
	// Constraints for the authenticator
	Constraints map[string]interface{} `json:"constraints,omitempty"`
	// Enrollment requirements for the authenticator
	Enroll map[string]interface{} `json:"enroll,omitempty"`
	// The authenticator ID for `custom_app`, `custom_otp` or `external_idp`. Use this property to select a specific `custom_app`, `custom_otp` or `external_idp` authenticator.
	ID string `json:"id,omitempty"`
	Key interface{} `json:"key,omitempty"`
}
