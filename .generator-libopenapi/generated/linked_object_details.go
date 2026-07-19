// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LinkedObjectDetails represents the LinkedObjectDetails schema
type LinkedObjectDetails struct {
	Type interface{} `json:"type"`
	// Description of the `primary` or the `associated` relationship
	Description string `json:"description,omitempty"`
	// API name of the `primary` or the `associated` link. The `name` parameter can't start with a number and can only contain the following characters: `a-z`, `A-Z`,` 0-9`, and `_`.
	Name string `json:"name"`
	// Display name of the `primary` or the `associated` link
	Title string `json:"title"`
}
