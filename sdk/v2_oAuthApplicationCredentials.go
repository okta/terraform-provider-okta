package sdk

type OAuthApplicationCredentials struct {
	Signing          *ApplicationCredentialsSigning          `json:"signing,omitempty"`
	UserNameTemplate *ApplicationCredentialsUsernameTemplate `json:"userNameTemplate,omitempty"`
	OauthClient      *ApplicationCredentialsOAuthClient      `json:"oauthClient,omitempty"`
}
