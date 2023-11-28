// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AccessPolicyRuleApplicationSignOn struct {
	Access             string              `json:"access,omitempty"`
	VerificationMethod *VerificationMethod `json:"verificationMethod,omitempty"`
}

func NewAccessPolicyRuleApplicationSignOn() *AccessPolicyRuleApplicationSignOn {
	return &AccessPolicyRuleApplicationSignOn{}
}

func (a *AccessPolicyRuleApplicationSignOn) IsPolicyInstance() bool {
	return true
}
