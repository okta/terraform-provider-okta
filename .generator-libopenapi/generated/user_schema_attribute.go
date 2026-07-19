// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserSchemaAttribute represents the UserSchemaAttribute schema
type UserSchemaAttribute struct {
	// Identifies the type of data represented by the string
	Format interface{} `json:"format,omitempty"`
	// Minimum character length of a string property
	MinLength int `json:"minLength,omitempty"`
	// Determines whether property values must be unique
	Unique bool `json:"unique,omitempty"`
	// Name of the property as it exists in an external application  **NOTE**: When you add a custom property, only Identity Provider app user schemas require `externalName` to be included in the request ...
	ExternalName string `json:"externalName,omitempty"`
	// Identifies where the property is mastered
	Master interface{} `json:"master,omitempty"`
	// Maximum character length of a string property
	MaxLength int `json:"maxLength,omitempty"`
	// Non-empty array of valid JSON schemas.  The `oneOf` key is only supported in conjunction with `enum` and provides a mechanism to return a display name for the `enum` value.<br> Each schema has the ...
	OneOf []interface{} `json:"oneOf,omitempty"`
	// Access control permissions for the property
	Permissions []interface{} `json:"permissions,omitempty"`
	Scope interface{} `json:"scope,omitempty"`
	// Enumerated value of the property.  The value of the property is limited to one of the values specified in the enum definition. The list of values for the enum must consist of unique elements.
	Enum []interface{} `json:"enum,omitempty"`
	// Type of property
	Type interface{} `json:"type,omitempty"`
	// If specified, assigns the value as the default value for the custom attribute. This is a nullable property. If you don't specify a value for this custom attribute during user creation or update, th...
	Default interface{} `json:"default,omitempty"`
	// Description of the property
	Description string `json:"description,omitempty"`
	// Namespace from the external application
	ExternalNamespace string `json:"externalNamespace,omitempty"`
	// Defines the mutability of the property
	Mutability interface{} `json:"mutability,omitempty"`
	// For `string` property types, specifies the regular expression used to validate the property
	Pattern string `json:"pattern,omitempty"`
	// Determines whether the property is required
	Required bool `json:"required,omitempty"`
	// User-defined display name for the property
	Title string `json:"title,omitempty"`
}
