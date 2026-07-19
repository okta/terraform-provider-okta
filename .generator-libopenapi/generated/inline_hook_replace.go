// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookReplace represents the InlineHookReplace schema
// An inline hook object that specifies the details of the inline hook
type InlineHookReplace struct {
	Channel interface{} `json:"channel,omitempty"`
	// The display name of the inline hook
	Name string `json:"name,omitempty"`
	// Version of the inline hook type. The currently supported version is `1.0.0`.
	Version string `json:"version,omitempty"`
}
