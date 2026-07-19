// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticatorEnrollmentPolicyRuleConditions represents the AuthenticatorEnrollmentPolicyRuleConditions schema
// Specifies conditions that must be met during policy evaluation to apply the rule. All policy conditions and conditions for at least one rule must be met to apply the settings specified in the polic...
type AuthenticatorEnrollmentPolicyRuleConditions struct {
	App interface{} `json:"app,omitempty"`
	Network interface{} `json:"network,omitempty"`
	// Identifies users and groups that are used together
	People map[string]interface{} `json:"people,omitempty"`
}
