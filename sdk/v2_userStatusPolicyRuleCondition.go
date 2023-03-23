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
