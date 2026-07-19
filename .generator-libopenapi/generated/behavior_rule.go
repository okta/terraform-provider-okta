// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BehaviorRule represents the BehaviorRule schema
type BehaviorRule struct {
	// ID of the Behavior Detection Rule
	ID string `json:"id,omitempty"`
	// Timestamp when the Behavior Detection Rule was last modified
	LastUpdated string `json:"lastUpdated,omitempty"`
	// Name of the Behavior Detection Rule
	Name string `json:"name"`
	Status interface{} `json:"status,omitempty"`
	Type interface{} `json:"type"`
	Link interface{} `json:"_link,omitempty"`
	// Timestamp when the Behavior Detection Rule was created
	Created string `json:"created,omitempty"`
}
