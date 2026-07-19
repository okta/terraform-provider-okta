// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CreateGroupRuleRequest represents the CreateGroupRuleRequest schema
type CreateGroupRuleRequest struct {
	Actions interface{} `json:"actions,omitempty"`
	Conditions interface{} `json:"conditions,omitempty"`
	// Name of the group rule
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}
