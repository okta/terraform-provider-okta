// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EntitlementValue represents the EntitlementValue schema
type EntitlementValue struct {
	// Entitlement value ID
	ID string `json:"id,omitempty"`
	// The entitlement value resource name
	Name string `json:"name,omitempty"`
	// The entitlement value resource [ORN](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn)
	Value string `json:"value,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json...
	Links map[string]interface{} `json:"_links,omitempty"`
}
