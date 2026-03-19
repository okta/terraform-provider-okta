// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AuthenticatorProvider struct {
	Configuration *AuthenticatorProviderConfiguration `json:"configuration,omitempty"`
	Type          string                              `json:"type,omitempty"`
}
