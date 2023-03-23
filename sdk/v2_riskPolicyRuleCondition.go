package sdk

type RiskPolicyRuleCondition struct {
	Behaviors []string `json:"behaviors,omitempty"`
}

func NewRiskPolicyRuleCondition() *RiskPolicyRuleCondition {
	return &RiskPolicyRuleCondition{}
}

func (a *RiskPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
