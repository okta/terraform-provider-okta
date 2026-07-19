// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RegistrationResponse represents the RegistrationResponse schema
type RegistrationResponse struct {
	// The `commands` object lets you invoke commands to modify or add values to the attributes in the Okta user profile that are created for this user. The object also lets you control whether or not the...
	Commands []map[string]interface{} `json:"commands,omitempty"`
	// For the registration inline hook, the `error` object provides a way of displaying an error message to the end user who is trying to register or update their profile.  * If you're using the Okta Sig...
	Error map[string]interface{} `json:"Error,omitempty"`
}
