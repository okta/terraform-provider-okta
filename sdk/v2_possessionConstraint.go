package sdk

type PossessionConstraint struct {
	Methods            []string `json:"methods,omitempty"`
	ReauthenticateIn   string   `json:"reauthenticateIn,omitempty"`
	Types              []string `json:"types,omitempty"`
	DeviceBound        string   `json:"deviceBound,omitempty"`
	HardwareProtection string   `json:"hardwareProtection,omitempty"`
	PhishingResistant  string   `json:"phishingResistant,omitempty"`
	UserPresence       string   `json:"userPresence,omitempty"`
}

func NewPossessionConstraint() *PossessionConstraint {
	return &PossessionConstraint{}
}

func (a *PossessionConstraint) IsPolicyInstance() bool {
	return true
}
