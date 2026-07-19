// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// Application represents the Application schema
type Application struct {
	// The Okta resource name (ORN) for the current app instance
	Orn string `json:"orn,omitempty"`
	UniversalLogout interface{} `json:"universalLogout,omitempty"`
	// Embedded resources related to the app using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification. If the `expand=user/{userId}` que...
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	ExpressConfiguration interface{} `json:"expressConfiguration,omitempty"`
	Licensing interface{} `json:"licensing,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Visibility interface{} `json:"visibility,omitempty"`
	// Timestamp when the application object was created
	Created *time.Time `json:"created,omitempty"`
	// Timestamp when the application object was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	SignOnMode interface{} `json:"signOnMode"`
	Links interface{} `json:"_links,omitempty"`
	// Enabled app features > **Note:** See [Application Features](/openapi/okta-management/management/tags/applicationfeatures/) for app provisioning features.
	Features []string `json:"features,omitempty"`
	// Unique ID for the app instance
	ID string `json:"id,omitempty"`
	// Contains any valid JSON schema for specifying properties that can be referenced from a request (only available to OAuth 2.0 client apps). For example, add an app manager contact email address or de...
	Profile map[string]interface{} `json:"profile,omitempty"`
	Accessibility interface{} `json:"accessibility,omitempty"`
	Label interface{} `json:"label"`
}
