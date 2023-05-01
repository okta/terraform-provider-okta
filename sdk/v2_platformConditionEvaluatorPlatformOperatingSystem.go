package sdk

type PlatformConditionEvaluatorPlatformOperatingSystem struct {
	Expression string                                                    `json:"expression,omitempty"`
	Type       string                                                    `json:"type,omitempty"`
	Version    *PlatformConditionEvaluatorPlatformOperatingSystemVersion `json:"version,omitempty"`
}

func NewPlatformConditionEvaluatorPlatformOperatingSystem() *PlatformConditionEvaluatorPlatformOperatingSystem {
	return &PlatformConditionEvaluatorPlatformOperatingSystem{}
}

func (a *PlatformConditionEvaluatorPlatformOperatingSystem) IsPolicyInstance() bool {
	return true
}
