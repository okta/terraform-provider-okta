package sdk

type ApplicationAccessibility struct {
	ErrorRedirectUrl string `json:"errorRedirectUrl,omitempty"`
	LoginRedirectUrl string `json:"loginRedirectUrl,omitempty"`
	SelfService      *bool  `json:"selfService,omitempty"`
}
