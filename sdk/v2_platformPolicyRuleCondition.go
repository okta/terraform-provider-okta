// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
