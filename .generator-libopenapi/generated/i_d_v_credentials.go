// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IDVCredentials represents the IDVCredentials schema
// Credentials for verifying requests to the IDV vendor
type IDVCredentials struct {
	// Client credential for `IDV_PERSONA` IdP type
	Bearer map[string]interface{} `json:"bearer,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Client credentials for `IDV_CLEAR` and `IDV_INCODE` IdP types
	Client map[string]interface{} `json:"client,omitempty"`
}
