package sdk

type UserActivationToken struct {
	ActivationToken string `json:"activationToken,omitempty"`
	ActivationUrl   string `json:"activationUrl,omitempty"`
}
