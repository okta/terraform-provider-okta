// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ResourceSetBindings represents the ResourceSetBindings schema
type ResourceSetBindings struct {
	// Roles associated with the resource set binding. If there are more than 100 bindings for the specified resource set, then the `_links.next` resource is returned with the next list of bindings.
	Roles []interface{} `json:"roles,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
