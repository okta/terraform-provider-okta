// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// RoleGovernanceSource represents the RoleGovernanceSource schema
// User role governance source
type RoleGovernanceSource struct {
	// The expiration date of the entitlement bundle
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
	// `id` of the grant
	GrantId string `json:"grantId"`
	Type interface{} `json:"type"`
	Links interface{} `json:"_links,omitempty"`
	// `id` of the entitlement bundle
	BundleId string `json:"bundleId,omitempty"`
}
