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
