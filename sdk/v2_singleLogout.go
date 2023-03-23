package sdk

type SingleLogout struct {
	Enabled   *bool  `json:"enabled,omitempty"`
	Issuer    string `json:"issuer,omitempty"`
	LogoutUrl string `json:"logoutUrl,omitempty"`
}
