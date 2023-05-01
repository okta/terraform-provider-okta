package sdk

type CreateUserRequest struct {
	Credentials *UserCredentials `json:"credentials,omitempty"`
	GroupIds    []string         `json:"groupIds,omitempty"`
	Profile     *UserProfile     `json:"profile,omitempty"`
	Type        *UserType        `json:"type,omitempty"`
}
