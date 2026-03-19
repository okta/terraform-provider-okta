// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ProfileEnrollmentPolicyRuleProfileAttribute struct {
	Label    string `json:"label,omitempty"`
	Name     string `json:"name,omitempty"`
	Required *bool  `json:"required,omitempty"`
}

func NewProfileEnrollmentPolicyRuleProfileAttribute() *ProfileEnrollmentPolicyRuleProfileAttribute {
	return &ProfileEnrollmentPolicyRuleProfileAttribute{}
}

func (a *ProfileEnrollmentPolicyRuleProfileAttribute) IsPolicyInstance() bool {
	return true
}
