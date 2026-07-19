// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TokenHookResponse represents the TokenHookResponse schema
// For the token inline hook, the `commands` and `error` objects that you can return in the JSON payload of your response are defined in the following sections. > **Note:** The size of your response p...
type TokenHookResponse struct {
	// You can use the `commands` object to provide commands to Okta. It's where you can tell Okta to add more claims to the token. The `commands` object is an array, allowing you to send multiple command...
	Commands []map[string]interface{} `json:"commands,omitempty"`
	// When an error object is returned, it causes Okta to return an OAuth 2.0 error to the requester of the token. In the error response, the value of `error` is `server_error`, and the value of `error_d...
	Error map[string]interface{} `json:"error,omitempty"`
}
