// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type MDMEnrollmentPolicyRuleCondition struct {
	BlockNonSafeAndroid *bool  `json:"blockNonSafeAndroid,omitempty"`
	Enrollment          string `json:"enrollment,omitempty"`
}

func NewMDMEnrollmentPolicyRuleCondition() *MDMEnrollmentPolicyRuleCondition {
	return &MDMEnrollmentPolicyRuleCondition{}
}

func (a *MDMEnrollmentPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
