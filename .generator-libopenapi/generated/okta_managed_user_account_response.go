// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// OktaManagedUserAccountResponse represents the OktaManagedUserAccountResponse schema
// An Okta managed user account representing a Universal Directory user managed as a service account
type OktaManagedUserAccountResponse struct {
	// Timestamp when the Okta managed user account was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// A list of IDs of the Okta groups who own the Okta managed user account
	OwnerGroupIds []string `json:"ownerGroupIds,omitempty"`
	// A list of IDs of the Okta users who own the Okta managed user account
	OwnerUserIds []string `json:"ownerUserIds,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// The username associated with the Okta user. This parameter is read-only, and it is derived from the Okta user profile.
	Username string `json:"username"`
	// Timestamp when the Okta managed user account was created
	Created *time.Time `json:"created,omitempty"`
	// The UUID of the Okta managed user account
	ID string `json:"id"`
	// The user-defined name for the Okta managed user account
	Name string `json:"name"`
	// The ID of the Okta user being managed as a service account
	OktaUserId string `json:"oktaUserId"`
	StatusDetail interface{} `json:"statusDetail,omitempty"`
	// The description of the Okta managed user account
	Description string `json:"description,omitempty"`
	// The email address associated with the Okta user. This parameter is read-only, and it is derived from the Okta user profile.
	Email string `json:"email"`
}
