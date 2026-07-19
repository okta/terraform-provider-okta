// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BulkGroupUpsertRequestBody represents the BulkGroupUpsertRequestBody schema
type BulkGroupUpsertRequestBody struct {
	// Array of group profiles that needs to be inserted or updated in Okta
	Profiles []map[string]interface{} `json:"profiles,omitempty"`
}
