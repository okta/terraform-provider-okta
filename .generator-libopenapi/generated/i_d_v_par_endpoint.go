// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IDVParEndpoint represents the IDVParEndpoint schema
// IDV [PAR](https://datatracker.ietf.org/doc/html/rfc9126) endpoint
type IDVParEndpoint struct {
	Binding string `json:"binding,omitempty"`
	// URL of the `par` endpoint of the IDV vendor
	Url string `json:"url,omitempty"`
}
