// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuthClientSecretLinks represents the OAuthClientSecretLinks schema
// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for the current status of an app using the [JSON Hypertext Application Language](https://datatracker.i...
type OAuthClientSecretLinks struct {
	Activate interface{} `json:"activate,omitempty"`
	Deactivate interface{} `json:"deactivate,omitempty"`
	Delete interface{} `json:"delete,omitempty"`
}
