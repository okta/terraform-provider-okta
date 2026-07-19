// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupQueryRequest represents the GroupQueryRequest schema
type GroupQueryRequest struct {
	// An array of LDAP group attribute names to retrieve. Restricted attributes: member, memberOf, *
	Attributes []string `json:"attributes"`
}
