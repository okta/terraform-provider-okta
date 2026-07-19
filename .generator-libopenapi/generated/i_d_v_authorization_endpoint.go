// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IDVAuthorizationEndpoint represents the IDVAuthorizationEndpoint schema
// IDV authorization endpoint
type IDVAuthorizationEndpoint struct {
	Binding string `json:"binding,omitempty"`
	// URL of the `authorization` endpoint of the IDV vendor
	Url string `json:"url,omitempty"`
}
