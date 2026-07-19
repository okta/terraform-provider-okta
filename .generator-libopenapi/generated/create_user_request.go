// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CreateUserRequest represents the CreateUserRequest schema
type CreateUserRequest struct {
	Credentials interface{} `json:"credentials,omitempty"`
	// The list of group IDs of groups that the user is added to at the time of creation
	GroupIds []string `json:"groupIds,omitempty"`
	Profile interface{} `json:"profile"`
	// The ID of the realm in which the user is residing. See [Realms](/openapi/okta-management/management/tags/realm).
	RealmId string `json:"realmId,omitempty"`
	// The ID of the user type. Add this value if you want to create a user with a non-default [User Type](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/UserType/). The user t...
	Type map[string]interface{} `json:"type,omitempty"`
}
