// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type DevicePolicyRuleConditionPlatform struct {
	SupportedMDMFrameworks []string `json:"supportedMDMFrameworks,omitempty"`
	Types                  []string `json:"types,omitempty"`
}

func NewDevicePolicyRuleConditionPlatform() *DevicePolicyRuleConditionPlatform {
	return &DevicePolicyRuleConditionPlatform{}
}

func (a *DevicePolicyRuleConditionPlatform) IsPolicyInstance() bool {
	return true
}
