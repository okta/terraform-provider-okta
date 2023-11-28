// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
