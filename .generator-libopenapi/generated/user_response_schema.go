// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UserResponseSchema represents the UserResponseSchema schema
type UserResponseSchema struct {
	Profile interface{} `json:"profile,omitempty"`
	// The timestamp when the user was created in the identity source
	Created *time.Time `json:"created,omitempty"`
	// The external ID of the user in the identity source
	ExternalId string `json:"externalId,omitempty"`
	// The ID of the user in the identity source
	ID string `json:"id,omitempty"`
	// The timestamp when the user was last updated in the identity source
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
}
