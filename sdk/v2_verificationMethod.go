package sdk

type VerificationMethod struct {
	Constraints      []*AccessPolicyConstraints `json:"constraints,omitempty"`
	FactorMode       string                     `json:"factorMode,omitempty"`
	InactivityPeriod string                     `json:"inactivityPeriod,omitempty"`
	ReauthenticateIn string                     `json:"reauthenticateIn,omitempty"`
	Type             string                     `json:"type,omitempty"`
}

func NewVerificationMethod() *VerificationMethod {
	return &VerificationMethod{}
}

func (a *VerificationMethod) IsPolicyInstance() bool {
	return true
}
