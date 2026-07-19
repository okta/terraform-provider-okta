// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupSchemaAttribute represents the GroupSchemaAttribute schema
type GroupSchemaAttribute struct {
	// Description of the property
	Description string `json:"description,omitempty"`
	// Name of the property as it exists in an external application
	ExternalName string `json:"externalName,omitempty"`
	// Type of property
	Type interface{} `json:"type,omitempty"`
	// Determines whether property values must be unique
	Unique bool `json:"unique,omitempty"`
	// Enumerated value of the property.  The value of the property is limited to one of the values specified in the enum definition. The list of values for the enum must consist of unique elements.
	Enum []interface{} `json:"enum,omitempty"`
	Items interface{} `json:"items,omitempty"`
	// Determines whether the property is required
	Required bool `json:"required,omitempty"`
	// Determines whether a group attribute can be set at the individual or group level
	Scope interface{} `json:"scope,omitempty"`
	// Access control permissions for the property
	Permissions []interface{} `json:"permissions,omitempty"`
	// Namespace from the external application
	ExternalNamespace string `json:"externalNamespace,omitempty"`
	// Identifies the type of data represented by the string
	Format interface{} `json:"format,omitempty"`
	// Identifies where the property is mastered
	Master interface{} `json:"master,omitempty"`
	// Maximum character length of a string property
	MaxLength int `json:"maxLength,omitempty"`
	// Minimum character length of a string property
	MinLength int `json:"minLength,omitempty"`
	// Defines the mutability of the property
	Mutability interface{} `json:"mutability,omitempty"`
	// User-defined display name for the property
	Title string `json:"title,omitempty"`
	// Non-empty array of valid JSON schemas.  The `oneOf` key is only supported in conjunction with `enum` and provides a mechanism to return a display name for the `enum` value.<br> Each schema has the ...
	OneOf []interface{} `json:"oneOf,omitempty"`
}
