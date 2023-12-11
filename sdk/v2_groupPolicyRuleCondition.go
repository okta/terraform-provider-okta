// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type GroupPolicyRuleCondition struct {
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

func NewGroupPolicyRuleCondition() *GroupPolicyRuleCondition {
	return &GroupPolicyRuleCondition{}
}

func (a *GroupPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
