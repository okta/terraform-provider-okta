// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserIdentifierPolicyRuleCondition represents the UserIdentifierPolicyRuleCondition schema
// Specifies a user identifier condition to match on
type UserIdentifierPolicyRuleCondition struct {
	// The name of the profile attribute to match against. Only used when type is `ATTRIBUTE`.
	Attribute string `json:"attribute,omitempty"`
	Patterns []interface{} `json:"patterns"`
	Type interface{} `json:"type"`
}
