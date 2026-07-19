// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GovernanceBundleCreateRequest represents the GovernanceBundleCreateRequest schema
// Request to create a governance bundle
type GovernanceBundleCreateRequest struct {
	// Description of the governance bundle
	Description string `json:"description,omitempty"`
	// List of entitlements to include in the governance bundle
	Entitlements []interface{} `json:"entitlements,omitempty"`
	// Name of the governance bundle
	Name string `json:"name,omitempty"`
}
