// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicyAccountLinkFilter represents the PolicyAccountLinkFilter schema
// Specifies filters on which users are available for account linking by an IdP
type PolicyAccountLinkFilter struct {
	Groups interface{} `json:"groups,omitempty"`
	Users interface{} `json:"users,omitempty"`
}
