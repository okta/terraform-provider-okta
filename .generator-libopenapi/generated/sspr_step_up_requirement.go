// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SsprStepUpRequirement represents the SsprStepUpRequirement schema
// Defines the secondary authenticators needed for password reset if `required` is true. The following are three valid configurations: * `required`=false * `required`=true with no methods to use any S...
type SsprStepUpRequirement struct {
	// Authenticator methods required for secondary authentication step of password recovery. Specify this value only when `required` is true and `security_question` is permitted for the secondary authent...
	Methods []string `json:"methods,omitempty"`
	Required bool `json:"required,omitempty"`
}
