// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BulkGroupDeleteRequestBody represents the BulkGroupDeleteRequestBody schema
type BulkGroupDeleteRequestBody struct {
	// Array of external IDs of groups that need to be deleted in Okta
	ExternalIds []string `json:"externalIds,omitempty"`
}
