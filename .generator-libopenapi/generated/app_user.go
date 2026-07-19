// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AppUser represents the AppUser schema
// The app user object defines a user's app-specific profile and credentials for an app
type AppUser struct {
	// Timestamp when the application user status was last changed
	StatusChanged *time.Time `json:"statusChanged,omitempty"`
	SyncState interface{} `json:"syncState,omitempty"`
	// Embedded resources related to the application user using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	Created interface{} `json:"created,omitempty"`
	Credentials interface{} `json:"credentials,omitempty"`
	// The ID of the user in the target app that's linked to the Okta application user object. This value is the native app-specific identifier or primary key for the user in the target app.  The `externa...
	ExternalId string `json:"externalId,omitempty"`
	// Unique identifier for the Okta user
	ID string `json:"id,omitempty"`
	// Timestamp when the application user password was last changed
	PasswordChanged *time.Time `json:"passwordChanged,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
	// Indicates if the assignment is direct (`USER`) or by group membership (`GROUP`). If not specified, Okta tries to determine the scope based on the assignment type.
	Scope string `json:"scope,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// Timestamp of the last synchronization operation. This value is only updated for apps with the `IMPORT_PROFILE_UPDATES` or `PUSH PROFILE_UPDATES` feature.
	LastSync *time.Time `json:"lastSync,omitempty"`
	LastUpdated interface{} `json:"lastUpdated,omitempty"`
}
