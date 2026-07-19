// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupsResponseSchema represents the GroupsResponseSchema schema
type GroupsResponseSchema struct {
	// The Okta group ID of the identity source group
	ID string `json:"id,omitempty"`
	// The profile information of the group
	Profile map[string]interface{} `json:"profile,omitempty"`
	// The external ID of the identity source group
	ExternalId string `json:"externalId,omitempty"`
}
