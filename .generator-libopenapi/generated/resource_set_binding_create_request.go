// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ResourceSetBindingCreateRequest represents the ResourceSetBindingCreateRequest schema
type ResourceSetBindingCreateRequest struct {
	// URLs to user and/or group instances that are assigned to the role
	Members []string `json:"members,omitempty"`
	// Unique key for the role
	Role string `json:"role,omitempty"`
}
