package sdk

type ContextPolicyRuleCondition struct {
	Expression string `json:"expression,omitempty"`
}

func NewContextPolicyRuleCondition() *ContextPolicyRuleCondition {
	return &ContextPolicyRuleCondition{}
}

func (a *ContextPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
