// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ApiToken represents the ApiToken schema
// An API token for an Okta User. This token is NOT scoped any further and can be used for any API the user has permissions to call.
type ApiToken struct {
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Name string `json:"name"`
	TokenWindow interface{} `json:"tokenWindow,omitempty"`
	UserId string `json:"userId,omitempty"`
	Link interface{} `json:"_link,omitempty"`
	ClientName string `json:"clientName,omitempty"`
	Created *time.Time `json:"created,omitempty"`
	ID string `json:"id,omitempty"`
	// The Network Condition of the API Token
	Network map[string]interface{} `json:"network,omitempty"`
}
