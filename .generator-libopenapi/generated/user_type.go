// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UserType represents the UserType schema
type UserType struct {
	Links interface{} `json:"_links,omitempty"`
	// A timestamp from when the user type was created
	Created *time.Time `json:"created,omitempty"`
	// The human-readable description of the user type
	Description string `json:"description,omitempty"`
	// A timestamp from when the user type was most recently updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The user ID of the most recent account to edit the user type
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	// The name of the user type. The name must start with A-Z or a-z and contain only A-Z, a-z, 0-9, or underscore (_) characters. This value becomes read-only after creation and can't be updated.
	Name string `json:"name"`
	// The user ID of the account that created the user type
	CreatedBy string `json:"createdBy,omitempty"`
	// A boolean value to indicate if this is the default user type
	Default bool `json:"default,omitempty"`
	// The human-readable name of the user type
	DisplayName string `json:"displayName"`
	// The unique key for the user type
	ID string `json:"id,omitempty"`
}
