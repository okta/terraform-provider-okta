// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// GroupRule represents the GroupRule schema
type GroupRule struct {
	Conditions interface{} `json:"conditions,omitempty"`
	// Creation date for group rule
	Created *time.Time `json:"created,omitempty"`
	// Date group rule was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Name of the group rule
	Name string `json:"name,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// Type to indicate a group rule operation. Only `group_rule` is allowed.
	Type string `json:"type,omitempty"`
	// This object appears with embedded resources related to the group rule if you use the `expand` query parameter
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Actions interface{} `json:"actions,omitempty"`
	// ID of the group rule
	ID string `json:"id,omitempty"`
}
