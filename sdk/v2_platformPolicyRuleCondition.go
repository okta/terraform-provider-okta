package sdk

type PlatformPolicyRuleCondition struct {
	Exclude []*PlatformConditionEvaluatorPlatform `json:"exclude,omitempty"`
	Include []*PlatformConditionEvaluatorPlatform `json:"include,omitempty"`
}

func NewPlatformPolicyRuleCondition() *PlatformPolicyRuleCondition {
	return &PlatformPolicyRuleCondition{}
}

func (a *PlatformPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
