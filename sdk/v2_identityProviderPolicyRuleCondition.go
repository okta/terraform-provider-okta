// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type IdentityProviderPolicyRuleCondition struct {
	IdpIds   []string `json:"idpIds,omitempty"`
	Provider string   `json:"provider,omitempty"`
}

func NewIdentityProviderPolicyRuleCondition() *IdentityProviderPolicyRuleCondition {
	return &IdentityProviderPolicyRuleCondition{}
}

func (a *IdentityProviderPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
