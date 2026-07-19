// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProfileMapping represents the ProfileMapping schema
// The profile mapping object describes a mapping between an Okta user's and an app user's properties using [JSON Schema Draft 4](https://datatracker.ietf.org/doc/html/draft-zyp-json-schema-04).  > **...
type ProfileMapping struct {
	// Unique identifier for a profile mapping
	ID string `json:"id,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Source interface{} `json:"source,omitempty"`
	Target interface{} `json:"target,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
