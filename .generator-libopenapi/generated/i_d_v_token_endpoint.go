// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IDVTokenEndpoint represents the IDVTokenEndpoint schema
// Token endpoint of the IDV vendor
type IDVTokenEndpoint struct {
	Binding string `json:"binding,omitempty"`
	// URL of the `token` endpoint of the IDV vendor
	Url string `json:"url,omitempty"`
}
