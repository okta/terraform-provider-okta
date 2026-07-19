// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AuthenticatorEnrollment represents the AuthenticatorEnrollment schema
type AuthenticatorEnrollment struct {
	// Status of the enrollment
	Status string `json:"status,omitempty"`
	Type interface{} `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Timestamp when the authenticator enrollment was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Timestamp when the authenticator enrollment was created
	Created *time.Time `json:"created,omitempty"`
	// The unique identifier of the authenticator enrollment
	ID string `json:"id,omitempty"`
	// A human-readable string that identifies the authenticator
	Key string `json:"key,omitempty"`
	// The authenticator display name
	Name string `json:"name,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
}
