// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// orgBillingContactType represents the orgBillingContactType schema
// Org billing contact
type orgBillingContactType struct {
	ContactType interface{} `json:"contactType,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for the org billing contact type object using the [JSON Hypertext Application Language](https://datatr...
	Links map[string]interface{} `json:"_links,omitempty"`
}
