// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicyAccountLinkFilterGroups represents the PolicyAccountLinkFilterGroups schema
// Group memberships used to determine link candidates
type PolicyAccountLinkFilterGroups struct {
	// Specifies the allowlist of Group identifiers to match against. Group memberships are restricted to type `OKTA_GROUP`.
	Include []string `json:"include,omitempty"`
}
