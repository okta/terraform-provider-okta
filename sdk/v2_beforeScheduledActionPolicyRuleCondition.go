// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type BeforeScheduledActionPolicyRuleCondition struct {
	Duration        *Duration                     `json:"duration,omitempty"`
	LifecycleAction *ScheduledUserLifecycleAction `json:"lifecycleAction,omitempty"`
}

func NewBeforeScheduledActionPolicyRuleCondition() *BeforeScheduledActionPolicyRuleCondition {
	return &BeforeScheduledActionPolicyRuleCondition{}
}

func (a *BeforeScheduledActionPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
