// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationLayout represents the ApplicationLayout schema
type ApplicationLayout struct {
	Label string `json:"label,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
	Rule map[string]interface{} `json:"rule,omitempty"`
	Scope string `json:"scope,omitempty"`
	Type string `json:"type,omitempty"`
	Elements []map[string]interface{} `json:"elements,omitempty"`
}
