// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UserFactorYubikeyOtpToken represents the UserFactorYubikeyOtpToken schema
type UserFactorYubikeyOtpToken struct {
	// Specified profile information for token
	Profile map[string]interface{} `json:"profile,omitempty"`
	// Token status
	Status string `json:"status,omitempty"`
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Timestamp when the token was created
	Created *time.Time `json:"created,omitempty"`
	// ID of the token
	ID string `json:"id,omitempty"`
	// Timestamp when the token was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Timestamp when the token was last verified
	LastVerified *time.Time `json:"lastVerified,omitempty"`
}
