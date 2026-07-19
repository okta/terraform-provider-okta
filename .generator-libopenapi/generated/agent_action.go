// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AgentAction represents the AgentAction schema
// Details about the Active Directory or LDAP group membership update
type AgentAction struct {
	// ID of the Active Directory or LDAP group to update
	ID string `json:"id"`
	Parameters interface{} `json:"parameters"`
}
