// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
