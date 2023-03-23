package sdk

type PasswordPolicyRuleAction struct {
	Access string `json:"access,omitempty"`
}

func NewPasswordPolicyRuleAction() *PasswordPolicyRuleAction {
	return &PasswordPolicyRuleAction{}
}

func (a *PasswordPolicyRuleAction) IsPolicyInstance() bool {
	return true
}
