// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaManagedUserAccountRequest represents the OktaManagedUserAccountRequest schema
// Request body for creating an Okta managed user account
type OktaManagedUserAccountRequest struct {
	// The description of the Okta managed user account
	Description string `json:"description,omitempty"`
	// The user-defined name for the Okta managed user account
	Name string `json:"name"`
	// The ID of the Okta user to manage as a service account. This must be an existing user in your Okta org.
	OktaUserId string `json:"oktaUserId"`
	// A list of IDs of the Okta groups who own the Okta managed user account
	OwnerGroupIds []string `json:"ownerGroupIds,omitempty"`
	// A list of IDs of the Okta users who own the Okta managed user account
	OwnerUserIds []string `json:"ownerUserIds,omitempty"`
}
