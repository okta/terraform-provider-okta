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
