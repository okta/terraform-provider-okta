// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdpPolicyRuleAction represents the IdpPolicyRuleAction schema
// Specifies where to route users when they are attempting to sign in to your org, if the rule conditions are satisfied. You can add up to 10 providers to a single `idp` policy action.
type IdpPolicyRuleAction struct {
	// Specifies IdP settings
	Idp map[string]interface{} `json:"idp,omitempty"`
}
