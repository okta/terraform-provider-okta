// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type VerifyFactorRequest struct {
	ActivationToken  string `json:"activationToken,omitempty"`
	Answer           string `json:"answer,omitempty"`
	Attestation      string `json:"attestation,omitempty"`
	ClientData       string `json:"clientData,omitempty"`
	NextPassCode     string `json:"nextPassCode,omitempty"`
	PassCode         string `json:"passCode,omitempty"`
	RegistrationData string `json:"registrationData,omitempty"`
	StateToken       string `json:"stateToken,omitempty"`
}

func NewVerifyFactorRequest() *VerifyFactorRequest {
	return &VerifyFactorRequest{}
}

func (a *VerifyFactorRequest) IsUserFactorInstance() bool {
	return true
}
