// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserImportResponse represents the UserImportResponse schema
type UserImportResponse struct {
	// The `commands` object is where you can provide commands to Okta. It is an array that allows you to send multiple commands. Each array element needs to consist of a type-value pair.
	Commands []map[string]interface{} `json:"commands,omitempty"`
	// An object to return an error. Returning an error causes Okta to record a failure event in the Okta System Log. The string supplied in the `errorSummary` property is recorded in the System Log event...
	Error map[string]interface{} `json:"error,omitempty"`
}
