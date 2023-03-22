package sdk

type AppUserCredentials struct {
	Password *AppUserPasswordCredential `json:"password,omitempty"`
	UserName string                     `json:"userName"`
}
