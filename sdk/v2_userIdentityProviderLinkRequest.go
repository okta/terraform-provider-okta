package sdk

type UserIdentityProviderLinkRequest struct {
	ExternalId string `json:"externalId,omitempty"`
}

func NewUserIdentityProviderLinkRequest() *UserIdentityProviderLinkRequest {
	return &UserIdentityProviderLinkRequest{}
}

func (a *UserIdentityProviderLinkRequest) IsPolicyInstance() bool {
	return true
}
