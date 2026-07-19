// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppServiceAccountForUpdate represents the AppServiceAccountForUpdate schema
type AppServiceAccountForUpdate struct {
	// The description of the app service account
	Description string `json:"description,omitempty"`
	// The user-defined name for the app service account
	Name string `json:"name,omitempty"`
	// A list of IDs of the Okta groups who own the app service account
	OwnerGroupIds []string `json:"ownerGroupIds,omitempty"`
	// A list of IDs of the Okta users who own the app service account
	OwnerUserIds []string `json:"ownerUserIds,omitempty"`
}
