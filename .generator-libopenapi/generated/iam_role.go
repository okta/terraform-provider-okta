// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// IamRole represents the IamRole schema
type IamRole struct {
	// Timestamp when the role was created
	Created *time.Time `json:"created,omitempty"`
	// Description of the role
	Description string `json:"description"`
	// Unique key for the role
	ID string `json:"id,omitempty"`
	// Unique label for the role
	Label string `json:"label"`
	// Timestamp when the role was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
