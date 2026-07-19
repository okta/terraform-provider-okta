// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// GroupOwner represents the GroupOwner schema
type GroupOwner struct {
	// Timestamp when the group owner was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The ID of the app instance if the `originType` is `APPLICATION`. This value is `NULL` if `originType` is `OKTA_DIRECTORY`.
	OriginId string `json:"originId,omitempty"`
	OriginType interface{} `json:"originType,omitempty"`
	// If `originType`is APPLICATION, this parameter is set to `FALSE` until the owner's `originId` is reconciled with an associated Okta ID.
	Resolved bool `json:"resolved,omitempty"`
	Type interface{} `json:"type,omitempty"`
	// The display name of the group owner
	DisplayName string `json:"displayName,omitempty"`
	// The `id` of the group owner
	ID string `json:"id,omitempty"`
}
