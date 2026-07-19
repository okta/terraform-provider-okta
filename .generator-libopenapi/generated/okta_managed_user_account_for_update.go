// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaManagedUserAccountForUpdate represents the OktaManagedUserAccountForUpdate schema
// Request body for updating an Okta managed user account
type OktaManagedUserAccountForUpdate struct {
	// A list of IDs of the Okta groups who own the Okta managed user account
	OwnerGroupIds []string `json:"ownerGroupIds,omitempty"`
	// A list of IDs of the Okta users who own the Okta managed user account
	OwnerUserIds []string `json:"ownerUserIds,omitempty"`
	// The description of the Okta managed user account
	Description string `json:"description,omitempty"`
	// The user-defined name for the Okta managed user account
	Name string `json:"name,omitempty"`
}
