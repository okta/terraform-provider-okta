// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationGroupAssignment represents the ApplicationGroupAssignment schema
// The Application Group object that defines a group of users' app-specific profile and credentials for an app
type ApplicationGroupAssignment struct {
	LastUpdated interface{} `json:"lastUpdated,omitempty"`
	// Priority assigned to the group. If an app has more than one group assigned to the same user, then the group with the higher priority has its profile applied to the [application user](https://develo...
	Priority int `json:"priority,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
	// Embedded resource related to the Application Group using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification. If the `expand=group...
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// ID of the [group](openapi/okta-management/management/group)
	ID string `json:"id,omitempty"`
}
