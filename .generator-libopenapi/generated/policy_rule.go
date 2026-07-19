// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// PolicyRule represents the PolicyRule schema
type PolicyRule struct {
	// Timestamp when the rule was last modified
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Priority of the rule
	Priority int `json:"priority,omitempty"`
	// Specifies whether Okta created the policy rule (`system=true`). You can't delete policy rules that have `system` set to `true`.
	System bool `json:"system,omitempty"`
	Type interface{} `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Identifier for the rule
	ID string `json:"id,omitempty"`
	// Name of the rule
	Name string `json:"name,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// Timestamp when the rule was created
	Created *time.Time `json:"created,omitempty"`
}
