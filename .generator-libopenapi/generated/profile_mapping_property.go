// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProfileMappingProperty represents the ProfileMappingProperty schema
// A target property, in string form, that maps to a valid [JSON Schema Draft](https://tools.ietf.org/html/draft-zyp-json-schema-04) document.
type ProfileMappingProperty struct {
	// Combination or single source properties that are mapped to the target property. See [Okta Expression Language](https://developer.okta.com/docs/reference/okta-expression-language/).
	Expression string `json:"expression,omitempty"`
	PushStatus interface{} `json:"pushStatus,omitempty"`
}
