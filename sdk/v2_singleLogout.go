// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type SingleLogout struct {
	Enabled   *bool  `json:"enabled,omitempty"`
	Issuer    string `json:"issuer,omitempty"`
	LogoutUrl string `json:"logoutUrl,omitempty"`
}
