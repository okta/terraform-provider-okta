// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
