// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PossessionConstraint struct {
	Methods            []string `json:"methods,omitempty"`
	ReauthenticateIn   string   `json:"reauthenticateIn,omitempty"`
	Types              []string `json:"types,omitempty"`
	DeviceBound        string   `json:"deviceBound,omitempty"`
	HardwareProtection string   `json:"hardwareProtection,omitempty"`
	PhishingResistant  string   `json:"phishingResistant,omitempty"`
	UserPresence       string   `json:"userPresence,omitempty"`
	UserVerification   string   `json:"userVerification,omitempty"`
}

func NewPossessionConstraint() *PossessionConstraint {
	return &PossessionConstraint{}
}

func (a *PossessionConstraint) IsPolicyInstance() bool {
	return true
}
