// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SessionIdentityProvider represents the SessionIdentityProvider schema
type SessionIdentityProvider struct {
	// Identity Provider ID. If the `type` is `OKTA`, then the `id` is the org ID.
	ID string `json:"id,omitempty"`
	Type interface{} `json:"type,omitempty"`
}
