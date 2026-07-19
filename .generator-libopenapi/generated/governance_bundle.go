// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GovernanceBundle represents the GovernanceBundle schema
type GovernanceBundle struct {
	// Description of the governance bundle
	Description string `json:"description,omitempty"`
	// Governance bundle ID
	ID string `json:"id,omitempty"`
	// Name of the governance bundle
	Name string `json:"name,omitempty"`
	// The governance bundle resource, in [ORN format](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn)
	Orn string `json:"orn,omitempty"`
	// Status of the governance bundle
	Status string `json:"status,omitempty"`
	// Link relations available
	Links map[string]interface{} `json:"_links,omitempty"`
}
