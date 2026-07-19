// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicyPeopleCondition represents the PolicyPeopleCondition schema
// Specifies the users and groups that are included or excluded by the policy rule
type PolicyPeopleCondition struct {
	Users interface{} `json:"users,omitempty"`
	Groups interface{} `json:"groups,omitempty"`
}
