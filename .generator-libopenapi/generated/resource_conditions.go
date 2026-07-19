// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ResourceConditions represents the ResourceConditions schema
// Conditions for further restricting a resource.
type ResourceConditions struct {
	// Specific resources to exclude
	Exclude map[string]interface{} `json:"Exclude,omitempty"`
}
