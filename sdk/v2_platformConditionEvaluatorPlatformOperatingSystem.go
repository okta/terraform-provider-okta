// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PlatformConditionEvaluatorPlatformOperatingSystem struct {
	Expression string                                                    `json:"expression"`
	Type       string                                                    `json:"type,omitempty"`
	Version    *PlatformConditionEvaluatorPlatformOperatingSystemVersion `json:"version,omitempty"`
}

func NewPlatformConditionEvaluatorPlatformOperatingSystem() *PlatformConditionEvaluatorPlatformOperatingSystem {
	return &PlatformConditionEvaluatorPlatformOperatingSystem{}
}

func (a *PlatformConditionEvaluatorPlatformOperatingSystem) IsPolicyInstance() bool {
	return true
}
