// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserPolicyRuleCondition represents the UserPolicyRuleCondition schema
// Specifies a set of Users to be included or excluded
type UserPolicyRuleCondition struct {
	LifecycleExpiration interface{} `json:"lifecycleExpiration,omitempty"`
	PasswordExpiration interface{} `json:"passwordExpiration,omitempty"`
	UserLifecycleAttribute interface{} `json:"userLifecycleAttribute,omitempty"`
	// Users to be excluded
	Exclude []string `json:"exclude,omitempty"`
	Inactivity interface{} `json:"inactivity,omitempty"`
	// Users to be included
	Include []string `json:"include,omitempty"`
}
