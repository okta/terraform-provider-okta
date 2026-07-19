// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogStreamSchema represents the LogStreamSchema schema
type LogStreamSchema struct {
	// URI of log stream schema
	ID string `json:"id,omitempty"`
	// Non-empty array of valid JSON schemas.  Okta only supports `oneOf` for specifying display names for an `enum`. Each schema has the following format:  ``` {   "const": "enumValue",   "title": "displ...
	OneOf []interface{} `json:"oneOf,omitempty"`
	// For `string` log stream schema property type, specifies the regular expression used to validate the property
	Pattern string `json:"pattern,omitempty"`
	// log stream schema properties object
	Properties map[string]interface{} `json:"properties,omitempty"`
	// Required properties for this log stream schema object
	Required []string `json:"required,omitempty"`
	// JSON schema version identifier
	$schema string `json:"$schema,omitempty"`
	// A collection of error messages for individual properties in the schema. Okta implements a subset of [ajv-errors](https://github.com/ajv-validator/ajv-errors).
	ErrorMessage map[string]interface{} `json:"errorMessage,omitempty"`
	// Name of the log streaming integration
	Title string `json:"title,omitempty"`
	// Type of log stream schema property
	Type string `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
