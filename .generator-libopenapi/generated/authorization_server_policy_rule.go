// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AuthorizationServerPolicyRule represents the AuthorizationServerPolicyRule schema
type AuthorizationServerPolicyRule struct {
	Conditions interface{} `json:"conditions,omitempty"`
	// Rule type
	Type string `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	Actions interface{} `json:"actions,omitempty"`
	// Timestamp when the rule was created
	Created *time.Time `json:"created,omitempty"`
	// Identifier of the rule
	ID string `json:"id,omitempty"`
	// Timestamp when the rule was last modified
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Name of the rule
	Name string `json:"name,omitempty"`
	// Priority of the rule
	Priority int `json:"priority,omitempty"`
	// Status of the rule
	Status string `json:"status,omitempty"`
	// Set to `true` for system rules. You can't delete system rules.
	System bool `json:"system,omitempty"`
}
