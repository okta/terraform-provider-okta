// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BundleEntitlementsResponse represents the BundleEntitlementsResponse schema
// Entitlement list for a governance bundle
type BundleEntitlementsResponse struct {
	// List of bundle entitlements
	Entitlements []interface{} `json:"entitlements,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
