// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// Realm represents the Realm schema
type Realm struct {
	Links interface{} `json:"_links,omitempty"`
	// Timestamp when the realm was created
	Created *time.Time `json:"created,omitempty"`
	// Unique ID for the realm
	ID string `json:"id,omitempty"`
	// Indicates the default realm. Existing users will start out in the default realm and can be moved to other realms individually or through realm assignments. See [Realms Assignments API](/openapi/okt...
	IsDefault bool `json:"isDefault,omitempty"`
	// Timestamp when the realm was updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
}
