// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdentityProviderPolicyRuleCondition represents the IdentityProviderPolicyRuleCondition schema
// Specifies the IdP that's used to sign in
type IdentityProviderPolicyRuleCondition struct {
	// Specifies the IdP ID
	IdpIds []string `json:"idpIds,omitempty"`
	Provider interface{} `json:"provider,omitempty"`
}
