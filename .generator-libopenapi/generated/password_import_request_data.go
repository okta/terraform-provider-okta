// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordImportRequestData represents the PasswordImportRequestData schema
type PasswordImportRequestData struct {
	// This object specifies the default action Okta is set to take. Okta takes this action if your external service sends an empty HTTP 204 response. You can override the default action by returning a co...
	Action map[string]interface{} `json:"action,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}
