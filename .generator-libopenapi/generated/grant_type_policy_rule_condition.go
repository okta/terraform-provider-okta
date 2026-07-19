// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GrantTypePolicyRuleCondition represents the GrantTypePolicyRuleCondition schema
// Array of grant types that this condition includes. Determines the mechanism that Okta uses to authorize the creation of the tokens.
type GrantTypePolicyRuleCondition struct {
	// Array of grant types that this condition includes.
	Include []string `json:"include,omitempty"`
}
