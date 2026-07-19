// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyRuleConditions represents the PasswordPolicyRuleConditions schema
// Specifies conditions that must be met during policy evaluation to apply the rule. All policy conditions and conditions for at least one rule must be met to apply the settings specified in the polic...
type PasswordPolicyRuleConditions struct {
	Network interface{} `json:"network,omitempty"`
	People interface{} `json:"people,omitempty"`
}
