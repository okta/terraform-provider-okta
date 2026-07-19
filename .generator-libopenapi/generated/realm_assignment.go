// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// RealmAssignment represents the RealmAssignment schema
type RealmAssignment struct {
	// Unique ID of the realm assignment
	ID string `json:"id,omitempty"`
	// Name of the realm
	Name string `json:"name,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	Actions interface{} `json:"actions,omitempty"`
	// Timestamp when the realm assignment was created
	Created *time.Time `json:"created,omitempty"`
	// Array of allowed domains. No user in this realm can be created or updated unless they have a username and email from one of these domains.  The following characters aren't allowed in the domain nam...
	Domains []string `json:"domains,omitempty"`
	// Indicates the default realm. Existing users will start out in the default realm and can be moved individually to other realms.
	IsDefault bool `json:"isDefault,omitempty"`
	// Timestamp of when the realm assignment was updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The priority of the realm assignment. The lower the number, the higher the priority. This helps resolve conflicts between realm assignments. > **Note:** When you create realm assignments in bulk, r...
	Priority int `json:"priority,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Conditions interface{} `json:"conditions,omitempty"`
}
