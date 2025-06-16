package sdk

type VerificationMethod struct {
	Constraints      []*AccessPolicyConstraints `json:"constraints,omitempty"`
	FactorMode       string                     `json:"factorMode,omitempty"`
	InactivityPeriod string                     `json:"inactivityPeriod,omitempty"`
	ReauthenticateIn string                     `json:"reauthenticateIn,omitempty"`
	Type             string                     `json:"type,omitempty"`
	ID               string                     `json:"id,omitempty"`
	Chains           []*AccessPolicyChains      `json:"chains,omitempty"`
}

type AccessPolicyChains struct {
	AuthenticationMethods []*AuthenticationMethodAccessPolicy `json:"authenticationMethods,omitempty"`
	Next                  []*AccessPolicyChains               `json:"next,omitempty"`
	ReauthenticateIn      string                              `json:"reauthenticateIn,omitempty"`
}

type AuthenticationMethodAccessPolicy struct {
	Key                     string   `json:"key,omitempty"`
	Method                  string   `json:"method,omitempty"`
	HardwareProtection      string   `json:"hardwareProtection,omitempty"`
	ID                      string   `json:"id,omitempty"`
	PhishingResistant       string   `json:"phishingResistant,omitempty"`
	UserVerification        string   `json:"userVerification,omitempty"`
	UserVerificationMethods []string `json:"userVerificationMethods,omitempty"`
}

func NewVerificationMethod() *VerificationMethod {
	return &VerificationMethod{}
}

func (a *VerificationMethod) IsPolicyInstance() bool {
	return true
}
