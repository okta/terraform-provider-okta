// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ClientPolicyCondition represents the ClientPolicyCondition schema
// Specifies which clients are included in the policy
type ClientPolicyCondition struct {
	// Which clients are included in the policy
	Include []string `json:"include,omitempty"`
}
