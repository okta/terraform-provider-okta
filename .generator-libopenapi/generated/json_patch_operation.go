// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// JsonPatchOperation represents the JsonPatchOperation schema
// The update action
type JsonPatchOperation struct {
	Op interface{} `json:"op,omitempty"`
	// The resource path of the attribute to update
	Path string `json:"path,omitempty"`
	// The update operation value
	Value map[string]interface{} `json:"value,omitempty"`
}
