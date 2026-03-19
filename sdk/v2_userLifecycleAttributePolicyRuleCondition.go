// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserLifecycleAttributePolicyRuleCondition struct {
	AttributeName string `json:"attributeName,omitempty"`
	MatchingValue string `json:"matchingValue,omitempty"`
}

func NewUserLifecycleAttributePolicyRuleCondition() *UserLifecycleAttributePolicyRuleCondition {
	return &UserLifecycleAttributePolicyRuleCondition{}
}

func (a *UserLifecycleAttributePolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
