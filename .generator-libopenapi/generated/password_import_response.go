// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordImportResponse represents the PasswordImportResponse schema
// Password import inline hook response
type PasswordImportResponse struct {
	// The `commands` object specifies whether Okta accepts the end user's sign-in credentials as valid or not. For the password import inline hook, you typically only return one `commands` object with on...
	Commands []map[string]interface{} `json:"commands,omitempty"`
}
