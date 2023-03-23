package sdk

type AppInstancePolicyRuleCondition struct {
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

func NewAppInstancePolicyRuleCondition() *AppInstancePolicyRuleCondition {
	return &AppInstancePolicyRuleCondition{}
}

func (a *AppInstancePolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
