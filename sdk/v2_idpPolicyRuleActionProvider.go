// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type IdpPolicyRuleActionProvider struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

func NewIdpPolicyRuleActionProvider() *IdpPolicyRuleActionProvider {
	return &IdpPolicyRuleActionProvider{}
}

func (a *IdpPolicyRuleActionProvider) IsPolicyInstance() bool {
	return true
}
