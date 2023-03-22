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
