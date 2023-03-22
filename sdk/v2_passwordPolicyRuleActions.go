package sdk

type PasswordPolicyRuleActions struct {
	Enroll                   *PolicyRuleActionsEnroll           `json:"enroll,omitempty"`
	Idp                      *IdpPolicyRuleAction               `json:"idp,omitempty"`
	PasswordChange           *PasswordPolicyRuleAction          `json:"passwordChange,omitempty"`
	SelfServicePasswordReset *PasswordPolicyRuleAction          `json:"selfServicePasswordReset,omitempty"`
	SelfServiceUnlock        *PasswordPolicyRuleAction          `json:"selfServiceUnlock,omitempty"`
	Signon                   *OktaSignOnPolicyRuleSignonActions `json:"signon,omitempty"`
}

func NewPasswordPolicyRuleActions() *PasswordPolicyRuleActions {
	return &PasswordPolicyRuleActions{}
}

func (a *PasswordPolicyRuleActions) IsPolicyInstance() bool {
	return true
}
