// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BundleEntitlement represents the BundleEntitlement schema
// An entitlement in a governance bundle
type BundleEntitlement struct {
	// The description of the role
	Description string `json:"description,omitempty"`
	// Entitlement ID
	ID string `json:"id,omitempty"`
	// The name of the role
	Name string `json:"name,omitempty"`
	// The role key
	Role string `json:"role,omitempty"`
	// Link relations available
	Links map[string]interface{} `json:"_links,omitempty"`
}
