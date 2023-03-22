package sdk

type UserLifecycleAttributePolicyRuleCondition struct {
	AttributeName string `json:"attributeName,omitempty"`
	MatchingValue string `json:"matchingValue,omitempty"`
}

func NewUserLifecycleAttributePolicyRuleCondition() *UserLifecycleAttributePolicyRuleCondition {
	return &UserLifecycleAttributePolicyRuleCondition{}
}

func (a *UserLifecycleAttributePolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
