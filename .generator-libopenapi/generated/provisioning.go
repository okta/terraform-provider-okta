// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Provisioning represents the Provisioning schema
// Specifies the behavior for just-in-time (JIT) provisioning of an IdP user as a new Okta user and their group memberships
type Provisioning struct {
	Conditions interface{} `json:"conditions,omitempty"`
	Groups interface{} `json:"groups,omitempty"`
	// Determines if the IdP should act as a source of truth for user profile attributes
	ProfileMaster bool `json:"profileMaster,omitempty"`
	Action interface{} `json:"action,omitempty"`
}
