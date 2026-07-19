// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CapabilitiesCreateObject represents the CapabilitiesCreateObject schema
// Determines whether Okta assigns a new app account to each user managed by Okta.  Okta doesn't create a new account if it detects that the username specified in Okta already exists in the app. The u...
type CapabilitiesCreateObject struct {
	LifecycleCreate interface{} `json:"lifecycleCreate,omitempty"`
}
