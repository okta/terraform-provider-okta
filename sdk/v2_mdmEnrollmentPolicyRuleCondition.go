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
