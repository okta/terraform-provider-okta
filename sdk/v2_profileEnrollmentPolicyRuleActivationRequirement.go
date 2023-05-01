package sdk

type ProfileEnrollmentPolicyRuleActivationRequirement struct {
	EmailVerification *bool `json:"emailVerification,omitempty"`
}

func NewProfileEnrollmentPolicyRuleActivationRequirement() *ProfileEnrollmentPolicyRuleActivationRequirement {
	return &ProfileEnrollmentPolicyRuleActivationRequirement{}
}

func (a *ProfileEnrollmentPolicyRuleActivationRequirement) IsPolicyInstance() bool {
	return true
}
