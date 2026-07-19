// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// TacAuthenticatorEnrollment represents the TacAuthenticatorEnrollment schema
type TacAuthenticatorEnrollment struct {
	// Timestamp when the authenticator enrollment was created
	Created *time.Time `json:"created,omitempty"`
	// A unique identifier of the authenticator enrollment
	ID string `json:"id,omitempty"`
	// A human-readable string that identifies the authenticator
	Key string `json:"key,omitempty"`
	// A user-friendly name for the authenticator enrollment
	Nickname string `json:"nickname,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
	// Status of the enrollment
	Status string `json:"status,omitempty"`
	Type interface{} `json:"type,omitempty"`
	// Timestamp when the authenticator enrollment was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The authenticator display name
	Name string `json:"name,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
