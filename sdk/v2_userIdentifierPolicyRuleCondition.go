package sdk

type UserIdentifierPolicyRuleCondition struct {
	Attribute string                                     `json:"attribute,omitempty"`
	Patterns  []*UserIdentifierConditionEvaluatorPattern `json:"patterns,omitempty"`
	Type      string                                     `json:"type,omitempty"`
}

func NewUserIdentifierPolicyRuleCondition() *UserIdentifierPolicyRuleCondition {
	return &UserIdentifierPolicyRuleCondition{}
}

func (a *UserIdentifierPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
