// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserStatusPolicyRuleCondition struct {
	Value string `json:"value,omitempty"`
}

func NewUserStatusPolicyRuleCondition() *UserStatusPolicyRuleCondition {
	return &UserStatusPolicyRuleCondition{}
}

func (a *UserStatusPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
