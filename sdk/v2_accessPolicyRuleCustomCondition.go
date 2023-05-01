package sdk

type AccessPolicyRuleCustomCondition struct {
	Condition string `json:"condition,omitempty"`
}

func NewAccessPolicyRuleCustomCondition() *AccessPolicyRuleCustomCondition {
	return &AccessPolicyRuleCustomCondition{}
}

func (a *AccessPolicyRuleCustomCondition) IsPolicyInstance() bool {
	return true
}
