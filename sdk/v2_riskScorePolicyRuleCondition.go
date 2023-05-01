package sdk

type RiskScorePolicyRuleCondition struct {
	Level string `json:"level,omitempty"`
}

func NewRiskScorePolicyRuleCondition() *RiskScorePolicyRuleCondition {
	return &RiskScorePolicyRuleCondition{}
}

func (a *RiskScorePolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
