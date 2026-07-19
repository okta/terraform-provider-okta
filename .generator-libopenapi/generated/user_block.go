// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserBlock represents the UserBlock schema
// Describes how the account is blocked from access. If `appliesTo` is `ANY_DEVICES`, then the account is blocked for all devices. If `appliesTo` is `UNKNOWN_DEVICES`, then the account is only blocked...
type UserBlock struct {
	// The devices that the block applies to
	AppliesTo string `json:"appliesTo,omitempty"`
	// Type of access block
	Type string `json:"type,omitempty"`
}
