// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// User represents the User schema
type User struct {
	// The timestamp when the user's password was last updated
	PasswordChanged *time.Time `json:"passwordChanged,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
	// The timestamp when the user status transitioned to `ACTIVE`
	Activated *time.Time `json:"activated,omitempty"`
	// The timestamp when the user was created
	Created *time.Time `json:"created,omitempty"`
	// Embedded resources related to the user using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	// The timestamp when the user was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The target status of an in-progress asynchronous status transition. This property is only returned if the user's state is transitioning.
	TransitioningToStatus string `json:"transitioningToStatus,omitempty"`
	// The user type that determines the schema for the user's profile. The `type` property is a map that identifies the [User Types](https://developer.okta.com/docs/api/openapi/okta-management/management...
	Type map[string]interface{} `json:"type,omitempty"`
	// Specifies link relations (see [Web Linking](https://datatracker.ietf.org/doc/html/rfc8288) available for the current status of a user. The links object is used for dynamic discovery of related reso...
	Links interface{} `json:"_links,omitempty"`
	Credentials interface{} `json:"credentials,omitempty"`
	// The ID of the realm in which the user is residing. See [Realms](/openapi/okta-management/management/tags/realm).
	RealmId string `json:"realmId,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// The timestamp when the status of the user last changed
	StatusChanged *time.Time `json:"statusChanged,omitempty"`
	// The unique key for the user
	ID string `json:"id,omitempty"`
	// The timestamp of the last login
	LastLogin *time.Time `json:"lastLogin,omitempty"`
}
