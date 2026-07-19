// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdpPolicyRuleActionProvider represents the IdpPolicyRuleActionProvider schema
type IdpPolicyRuleActionProvider struct {
	// IdP types of `OKTA`, `AgentlessDSSO`, and `IWA` don't require an ID.
	ID string `json:"id,omitempty"`
	// Provider `name` in Okta. Optional. Supported in `IDENTITY ENGINE`.
	Name string `json:"name,omitempty"`
	Type interface{} `json:"type,omitempty"`
}
