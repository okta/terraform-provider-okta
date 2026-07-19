// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IAMBundleEntitlement represents the IAMBundleEntitlement schema
// An entitlement in a governance bundle
type IAMBundleEntitlement struct {
	// List of resource set IDs for the custom role
	ResourceSets []string `json:"resourceSets,omitempty"`
	// The role
	Role string `json:"role,omitempty"`
	// List of target resource IDs to scope the entitlement with the role
	Targets []string `json:"targets,omitempty"`
}
