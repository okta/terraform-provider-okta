package sdk

type DeviceAccessPolicyRuleCondition struct {
	Assurance  *DeviceAssurancePolicyRuleCondition `json:"assurance,omitempty"`
	Migrated   *bool                               `json:"migrated,omitempty"`
	Platform   *DevicePolicyRuleConditionPlatform  `json:"platform,omitempty"`
	Rooted     *bool                               `json:"rooted,omitempty"`
	TrustLevel string                              `json:"trustLevel,omitempty"`
	Managed    *bool                               `json:"managed,omitempty"`
	Registered *bool                               `json:"registered,omitempty"`
}

func NewDeviceAccessPolicyRuleCondition() *DeviceAccessPolicyRuleCondition {
	return &DeviceAccessPolicyRuleCondition{}
}

func (a *DeviceAccessPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
