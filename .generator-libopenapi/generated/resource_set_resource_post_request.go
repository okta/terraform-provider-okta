// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ResourceSetResourcePostRequest represents the ResourceSetResourcePostRequest schema
type ResourceSetResourcePostRequest struct {
	Conditions interface{} `json:"conditions"`
	// Resource in ORN or REST API URL format
	ResourceOrnOrUrl string `json:"resourceOrnOrUrl"`
}
