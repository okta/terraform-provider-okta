package sdk

type AccessPolicyRuleApplicationSignOn struct {
	Access             string              `json:"access,omitempty"`
	VerificationMethod *VerificationMethod `json:"verificationMethod,omitempty"`
}

func NewAccessPolicyRuleApplicationSignOn() *AccessPolicyRuleApplicationSignOn {
	return &AccessPolicyRuleApplicationSignOn{}
}

func (a *AccessPolicyRuleApplicationSignOn) IsPolicyInstance() bool {
	return true
}
