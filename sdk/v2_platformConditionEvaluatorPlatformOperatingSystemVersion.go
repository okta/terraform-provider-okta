// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PlatformConditionEvaluatorPlatformOperatingSystemVersion struct {
	MatchType string `json:"matchType,omitempty"`
	Value     string `json:"value,omitempty"`
}

func NewPlatformConditionEvaluatorPlatformOperatingSystemVersion() *PlatformConditionEvaluatorPlatformOperatingSystemVersion {
	return &PlatformConditionEvaluatorPlatformOperatingSystemVersion{}
}

func (a *PlatformConditionEvaluatorPlatformOperatingSystemVersion) IsPolicyInstance() bool {
	return true
}
