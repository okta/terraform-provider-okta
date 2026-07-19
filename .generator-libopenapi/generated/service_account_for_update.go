// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ServiceAccountForUpdate represents the ServiceAccountForUpdate schema
type ServiceAccountForUpdate struct {
	// The human-readable name for the service account
	Name string `json:"name,omitempty"`
	// A list of IDs of the Okta groups who own the service account
	OwnerGroupIds []string `json:"ownerGroupIds,omitempty"`
	// A list of IDs of the Okta users who own the service account
	OwnerUserIds []string `json:"ownerUserIds,omitempty"`
	// The description of the service account
	Description string `json:"description,omitempty"`
}
