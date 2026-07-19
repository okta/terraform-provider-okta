// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GovernanceBundleUpdateRequest represents the GovernanceBundleUpdateRequest schema
// Request to update a governance bundle
type GovernanceBundleUpdateRequest struct {
	// Description of the governance bundle
	Description string `json:"description,omitempty"`
	// List of entitlements to include in the governance bundle
	Entitlements []interface{} `json:"entitlements,omitempty"`
	// Name of the governance bundle
	Name string `json:"name,omitempty"`
}
