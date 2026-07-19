// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppUserProfileRequestPayload represents the AppUserProfileRequestPayload schema
// Updates the assigned user profile > **Note:** The Okta API currently doesn't support entity tags for conditional updates. As long as you're the only user updating the the user profile, Okta recomme...
type AppUserProfileRequestPayload struct {
	Profile interface{} `json:"profile,omitempty"`
}
