// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type IdpPolicyRuleAction struct {
	Providers []*IdpPolicyRuleActionProvider `json:"providers,omitempty"`
}

func NewIdpPolicyRuleAction() *IdpPolicyRuleAction {
	return &IdpPolicyRuleAction{}
}

func (a *IdpPolicyRuleAction) IsPolicyInstance() bool {
	return true
}
