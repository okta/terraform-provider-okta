// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// DetailedHookKeyInstance represents the DetailedHookKeyInstance schema
// A key object with public key details
type DetailedHookKeyInstance struct {
	Embedded interface{} `json:"_embedded,omitempty"`
	// Timestamp when the key was created
	Created *time.Time `json:"created,omitempty"`
	// The unique Okta ID of this key record
	ID string `json:"id,omitempty"`
	// Whether this key is currently in use by other applications
	IsUsed string `json:"isUsed,omitempty"`
	// The alias of the public key
	KeyId string `json:"keyId,omitempty"`
	// Timestamp when the key was updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Display name of the key
	Name string `json:"name,omitempty"`
}
