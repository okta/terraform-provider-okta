// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SAMLHookResponse represents the SAMLHookResponse schema
type SAMLHookResponse struct {
	// The `commands` object is where you tell Okta to add additional claims to the assertion or to modify the existing assertion statements.  `commands` is an array, allowing you to send multiple command...
	Commands []map[string]interface{} `json:"commands,omitempty"`
	// An object to return an error. Returning an error causes Okta to record a failure event in the Okta System Log. The string supplied in the `errorSummary` property is recorded in the System Log event...
	Error map[string]interface{} `json:"error,omitempty"`
}
