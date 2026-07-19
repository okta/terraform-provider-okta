// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RealmProfile represents the RealmProfile schema
type RealmProfile struct {
	// Used to store partner users. This property must be set to `PARTNER` to access Okta's external partner portal.
	RealmType string `json:"realmType,omitempty"`
	// Array of allowed domains. No user in this realm can be created or updated unless they have a username and email from one of these domains.  The following characters aren't allowed in the domain nam...
	Domains []interface{} `json:"domains,omitempty"`
	// Name of a realm
	Name string `json:"name"`
}
