package sdk

type OpenIdConnectApplicationIdpInitiatedLogin struct {
	DefaultScope []string `json:"default_scope"`
	Mode         string   `json:"mode,omitempty"`
}
