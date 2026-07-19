// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// CreateOrUpdatePolicy represents the CreateOrUpdatePolicy schema
type CreateOrUpdatePolicy struct {
	// Timestamp when the policy was last modified
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Type interface{} `json:"type"`
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Timestamp when the policy was created
	Created *time.Time `json:"created,omitempty"`
	// Description of the policy
	Description string `json:"description,omitempty"`
	// Name of the policy
	Name string `json:"name"`
	// Specifies the order in which this policy is evaluated in relation to the other policies
	Priority int `json:"priority,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// Specifies whether Okta created the policy
	System bool `json:"system,omitempty"`
	// Identifier of the policy
	ID string `json:"id,omitempty"`
}
