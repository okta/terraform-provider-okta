// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TelephonyResponse represents the TelephonyResponse schema
// Telephony inline hook response
type TelephonyResponse struct {
	// The `commands` object specifies whether Okta accepts the end user's sign-in credentials as valid or not. For the telephony inline hook, you typically only return one `commands` object with one arra...
	Commands []map[string]interface{} `json:"commands,omitempty"`
}
