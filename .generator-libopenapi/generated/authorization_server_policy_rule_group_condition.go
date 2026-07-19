// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthorizationServerPolicyRuleGroupCondition represents the AuthorizationServerPolicyRuleGroupCondition schema
// Specifies a set of Groups whose Users are to be included
type AuthorizationServerPolicyRuleGroupCondition struct {
	// Groups to be included
	Include []string `json:"include,omitempty"`
}
