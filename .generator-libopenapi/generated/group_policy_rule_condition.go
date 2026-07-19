// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupPolicyRuleCondition represents the GroupPolicyRuleCondition schema
// Specifies a set of groups whose users are to be included or excluded
type GroupPolicyRuleCondition struct {
	// Groups to be excluded
	Exclude []string `json:"exclude,omitempty"`
	// Groups to be included
	Include []string `json:"include,omitempty"`
}
