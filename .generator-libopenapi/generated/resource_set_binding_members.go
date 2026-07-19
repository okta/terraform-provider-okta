// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ResourceSetBindingMembers represents the ResourceSetBindingMembers schema
type ResourceSetBindingMembers struct {
	// The members of the role resource set binding. If there are more than 100 members for the binding, then the `_links.next` resource is returned with the next list of members.
	Members []interface{} `json:"members,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
