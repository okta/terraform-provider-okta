// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type OAuthApplicationCredentials struct {
	Signing          *ApplicationCredentialsSigning          `json:"signing,omitempty"`
	UserNameTemplate *ApplicationCredentialsUsernameTemplate `json:"userNameTemplate,omitempty"`
	OauthClient      *ApplicationCredentialsOAuthClient      `json:"oauthClient,omitempty"`
}
