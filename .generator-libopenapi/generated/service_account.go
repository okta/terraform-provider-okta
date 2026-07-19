// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ServiceAccount represents the ServiceAccount schema
type ServiceAccount struct {
	Status interface{} `json:"status,omitempty"`
	// Timestamp when the service account was created
	Created *time.Time `json:"created,omitempty"`
	// The UUID of the service account
	ID string `json:"id,omitempty"`
	StatusDetail interface{} `json:"statusDetail,omitempty"`
	AccountType interface{} `json:"accountType"`
	// The description of the service account
	Description string `json:"description,omitempty"`
	// Timestamp when the service account was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The user-defined name for the service account
	Name string `json:"name"`
	// A list of IDs of the Okta groups that own the service account
	OwnerGroupIds []string `json:"ownerGroupIds,omitempty"`
	// A list of IDs of the Okta users that own the service account
	OwnerUserIds []string `json:"ownerUserIds,omitempty"`
}
