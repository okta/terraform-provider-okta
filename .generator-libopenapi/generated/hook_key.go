// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// HookKey represents the HookKey schema
// The `id` property in the response as `id` serves as the unique ID for the key, which you can specify when invoking other CRUD operations.  The `keyId` provided in the response is the alias of the p...
type HookKey struct {
	// Timestamp when the key was created
	Created *time.Time `json:"created,omitempty"`
	// The unique identifier for the key
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
