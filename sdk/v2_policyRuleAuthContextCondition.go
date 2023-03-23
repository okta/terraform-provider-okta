package sdk

type PolicyRuleAuthContextCondition struct {
	AuthType string `json:"authType,omitempty"`
}

func NewPolicyRuleAuthContextCondition() *PolicyRuleAuthContextCondition {
	return &PolicyRuleAuthContextCondition{}
}

func (a *PolicyRuleAuthContextCondition) IsPolicyInstance() bool {
	return true
}
