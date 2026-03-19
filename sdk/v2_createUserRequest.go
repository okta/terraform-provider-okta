// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type CreateUserRequest struct {
	Credentials *UserCredentials `json:"credentials,omitempty"`
	GroupIds    []string         `json:"groupIds,omitempty"`
	Profile     *UserProfile     `json:"profile,omitempty"`
	Type        *UserType        `json:"type,omitempty"`
	RealmId     *string          `json:"realmId,omitempty"`
}
