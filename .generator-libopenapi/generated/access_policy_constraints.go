// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AccessPolicyConstraints represents the AccessPolicyConstraints schema
// Specifies constraints for the authenticator. Constraints are logically evaluated such that only one constraint object needs to be satisfied. But, within a constraint object, each constraint propert...
type AccessPolicyConstraints struct {
	Knowledge interface{} `json:"knowledge,omitempty"`
	Possession interface{} `json:"possession,omitempty"`
}
