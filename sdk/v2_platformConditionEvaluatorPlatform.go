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
