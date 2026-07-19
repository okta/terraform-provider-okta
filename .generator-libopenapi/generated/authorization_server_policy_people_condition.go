// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthorizationServerPolicyPeopleCondition represents the AuthorizationServerPolicyPeopleCondition schema
// Identifies Users and Groups that are used together
type AuthorizationServerPolicyPeopleCondition struct {
	Groups interface{} `json:"groups,omitempty"`
	Users interface{} `json:"users,omitempty"`
}
