// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// RoleTarget represents the RoleTarget schema
type RoleTarget struct {
	Links interface{} `json:"_links,omitempty"`
	// The assignment type of how the user receives this target
	AssignmentType string `json:"assignmentType,omitempty"`
	// The expiry time stamp of the associated target. It's only included in the response if the associated target will expire.
	Expiration *time.Time `json:"expiration,omitempty"`
	// The [Okta Resource Name (ORN)](https://support.okta.com/help/s/article/understanding-okta-resource-name-orn) of the app target or group target
	Orn string `json:"orn,omitempty"`
}
