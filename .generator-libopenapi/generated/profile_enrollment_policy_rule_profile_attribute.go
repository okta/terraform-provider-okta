// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProfileEnrollmentPolicyRuleProfileAttribute represents the ProfileEnrollmentPolicyRuleProfileAttribute schema
type ProfileEnrollmentPolicyRuleProfileAttribute struct {
	// A display-friendly label for this property
	Label string `json:"label,omitempty"`
	// The name of a user profile property. Can be an existing property.
	Name string `json:"name,omitempty"`
	// (Optional, default `FALSE`) Indicates if this property is required for enrollment
	Required bool `json:"required,omitempty"`
}
