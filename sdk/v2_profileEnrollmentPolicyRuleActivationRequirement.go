// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
