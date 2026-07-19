// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UpdateUserRequest represents the UpdateUserRequest schema
type UpdateUserRequest struct {
	Credentials interface{} `json:"credentials,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
	// The ID of the realm in which the user is residing. See [Realms](/openapi/okta-management/management/tags/realm).
	RealmId string `json:"realmId,omitempty"`
	// The ID of the user type. Add this value if you want to create a user with a non-default [User Type](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/UserType/). The user t...
	Type map[string]interface{} `json:"type,omitempty"`
}
