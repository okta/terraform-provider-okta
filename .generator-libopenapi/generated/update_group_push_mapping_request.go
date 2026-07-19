// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UpdateGroupPushMappingRequest represents the UpdateGroupPushMappingRequest schema
type UpdateGroupPushMappingRequest struct {
	// The status of the group push mapping.  If changing the group push mapping status to `ACTIVE`, Okta performs an initial push to the target group, and then begins pushing membership changes.  If chan...
	Status string `json:"status"`
}
