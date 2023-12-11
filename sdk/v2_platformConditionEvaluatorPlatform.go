// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PlatformConditionEvaluatorPlatform struct {
	Os   *PlatformConditionEvaluatorPlatformOperatingSystem `json:"os,omitempty"`
	Type string                                             `json:"type,omitempty"`
}

func NewPlatformConditionEvaluatorPlatform() *PlatformConditionEvaluatorPlatform {
	return &PlatformConditionEvaluatorPlatform{}
}

func (a *PlatformConditionEvaluatorPlatform) IsPolicyInstance() bool {
	return true
}
