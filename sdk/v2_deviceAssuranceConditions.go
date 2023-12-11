// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type DeviceAssurancePolicyRuleCondition struct {
	Include []string `json:"include,omitempty"`
}

func NewDeviceAssurancePolicyRuleCondition() *DeviceAssurancePolicyRuleCondition {
	return &DeviceAssurancePolicyRuleCondition{}
}

func (a *DeviceAssurancePolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
