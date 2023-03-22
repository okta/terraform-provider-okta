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
