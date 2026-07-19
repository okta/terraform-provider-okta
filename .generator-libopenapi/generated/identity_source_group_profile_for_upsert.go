// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdentitySourceGroupProfileForUpsert represents the IdentitySourceGroupProfileForUpsert schema
// Contains a set of external group attributes and their values that are mapped to Okta standard properties. See the group [`profile` object](https://developer.okta.com/docs/api/openapi/okta-managemen...
type IdentitySourceGroupProfileForUpsert struct {
	// Description of the group
	Description string `json:"description,omitempty"`
	// Name of the group
	DisplayName string `json:"displayName,omitempty"`
}
