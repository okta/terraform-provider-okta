// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AuthenticatorBase represents the AuthenticatorBase schema
type AuthenticatorBase struct {
	// Timestamp when the authenticator was created
	Created *time.Time `json:"created,omitempty"`
	// A unique identifier for the authenticator
	ID string `json:"id,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="ea"></x-lifecycle></x-lifecycle-container>The description of the authenticator. This setting is only available for the `webauthn` authenticator type (Pass...
	Description string `json:"description,omitempty"`
	Key interface{} `json:"key,omitempty"`
	// Timestamp when the authenticator was last modified
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Display name of the authenticator
	Name string `json:"name,omitempty"`
	Type interface{} `json:"type,omitempty"`
	// Link relations for this object
	Links interface{} `json:"_links,omitempty"`
}
