// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// Group represents the Group schema
type Group struct {
	// Timestamp when the group's profile was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Determines the group's `profile`
	ObjectClass []string `json:"objectClass,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
	Type interface{} `json:"type,omitempty"`
	// Embedded resources related to the group
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	// Unique ID for the group
	ID string `json:"id,omitempty"`
	// Timestamp when the groups memberships were last updated
	LastMembershipUpdated *time.Time `json:"lastMembershipUpdated,omitempty"`
	// [Discoverable resources](/openapi/okta-management/management/group/listgroups#group/listgroups/t=response&c=200&path=_links/source) related to the group
	Links interface{} `json:"_links,omitempty"`
	// Timestamp when the group was created
	Created *time.Time `json:"created,omitempty"`
}
