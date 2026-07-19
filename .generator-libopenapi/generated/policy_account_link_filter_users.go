// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicyAccountLinkFilterUsers represents the PolicyAccountLinkFilterUsers schema
// Filters on which users are available for account linking
type PolicyAccountLinkFilterUsers struct {
	// Specifies the blocklist of user identifiers to exclude from account linking
	Exclude []string `json:"exclude,omitempty"`
	// Specifies whether admin users should be excluded from account linking
	ExcludeAdmins bool `json:"excludeAdmins,omitempty"`
}
