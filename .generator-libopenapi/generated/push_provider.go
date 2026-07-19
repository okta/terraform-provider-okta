// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PushProvider represents the PushProvider schema
type PushProvider struct {
	ProviderType interface{} `json:"providerType,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Unique key for the Push Provider
	ID string `json:"id,omitempty"`
	// Timestamp when the Push Provider was last modified
	LastUpdatedDate string `json:"lastUpdatedDate,omitempty"`
	// Display name of the push provider
	Name string `json:"name,omitempty"`
}
