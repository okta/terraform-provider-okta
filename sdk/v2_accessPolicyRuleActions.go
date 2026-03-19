// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AccessPolicyRuleActions struct {
	Enroll                   *PolicyRuleActionsEnroll           `json:"enroll,omitempty"`
	Idp                      *IdpPolicyRuleAction               `json:"idp,omitempty"`
	PasswordChange           *PasswordPolicyRuleAction          `json:"passwordChange,omitempty"`
	SelfServicePasswordReset *PasswordPolicyRuleAction          `json:"selfServicePasswordReset,omitempty"`
	SelfServiceUnlock        *PasswordPolicyRuleAction          `json:"selfServiceUnlock,omitempty"`
	Signon                   *OktaSignOnPolicyRuleSignonActions `json:"signon,omitempty"`
	AppSignOn                *AccessPolicyRuleApplicationSignOn `json:"appSignOn,omitempty"`
}

func NewAccessPolicyRuleActions() *AccessPolicyRuleActions {
	return &AccessPolicyRuleActions{}
}

func (a *AccessPolicyRuleActions) IsPolicyInstance() bool {
	return true
}
