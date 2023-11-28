// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AppAndInstancePolicyRuleCondition struct {
	Exclude []*AppAndInstanceConditionEvaluatorAppOrInstance `json:"exclude,omitempty"`
	Include []*AppAndInstanceConditionEvaluatorAppOrInstance `json:"include,omitempty"`
}

func NewAppAndInstancePolicyRuleCondition() *AppAndInstancePolicyRuleCondition {
	return &AppAndInstancePolicyRuleCondition{}
}

func (a *AppAndInstancePolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
