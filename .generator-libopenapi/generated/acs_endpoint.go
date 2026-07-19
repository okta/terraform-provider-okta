// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AcsEndpoint represents the AcsEndpoint schema
// An array of ACS endpoints. You can configure a maximum of 100 endpoints.
type AcsEndpoint struct {
	// Index of the URL in the array of ACS endpoints
	Index int `json:"index"`
	// URL of the ACS
	Url string `json:"url"`
}
