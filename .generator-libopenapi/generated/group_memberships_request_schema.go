// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupMembershipsRequestSchema represents the GroupMembershipsRequestSchema schema
type GroupMembershipsRequestSchema struct {
	// A list of app user external IDs to be inserted in this group in Okta
	MemberExternalIds []string `json:"memberExternalIds,omitempty"`
}
