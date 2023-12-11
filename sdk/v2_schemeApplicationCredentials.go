// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type SchemeApplicationCredentials struct {
	Signing          *ApplicationCredentialsSigning          `json:"signing,omitempty"`
	UserNameTemplate *ApplicationCredentialsUsernameTemplate `json:"userNameTemplate,omitempty"`
	Password         *PasswordCredential                     `json:"password,omitempty"`
	RevealPassword   *bool                                   `json:"revealPassword,omitempty"`
	Scheme           string                                  `json:"scheme,omitempty"`
	UserName         string                                  `json:"userName"`
}
