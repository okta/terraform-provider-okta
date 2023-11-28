// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type IdentityProviderCredentials struct {
	Client  *IdentityProviderCredentialsClient  `json:"client,omitempty"`
	Signing *IdentityProviderCredentialsSigning `json:"signing,omitempty"`
	Trust   *IdentityProviderCredentialsTrust   `json:"trust,omitempty"`
}
