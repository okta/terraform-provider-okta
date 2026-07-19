// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SsprPrimaryRequirement represents the SsprPrimaryRequirement schema
// Defines the authenticators permitted for the initial authentication step of password recovery
type SsprPrimaryRequirement struct {
	// Constraints on the values specified in the `methods` array. Specifying a constraint limits methods to specific authenticator(s). Currently, Google OTP is the only accepted constraint.
	MethodConstraints []interface{} `json:"methodConstraints,omitempty"`
	// Authenticator methods allowed for the initial authentication step of password recovery. Method `otp` requires a constraint limiting it to a Google authenticator.
	Methods []string `json:"methods,omitempty"`
}
