// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProvisioningGroups represents the ProvisioningGroups schema
// Provisioning settings for a user's group memberships
type ProvisioningGroups struct {
	// List of `OKTA_GROUP` group identifiers to add an IdP user as a member with the `ASSIGN` action
	Assignments []string `json:"assignments,omitempty"`
	// Allowlist of `OKTA_GROUP` group identifiers for the `APPEND` or `SYNC` provisioning action
	Filter []string `json:"filter,omitempty"`
	// IdP user profile attribute name (case-insensitive) for an array value that contains group memberships
	SourceAttributeName string `json:"sourceAttributeName,omitempty"`
	Action interface{} `json:"action,omitempty"`
}
