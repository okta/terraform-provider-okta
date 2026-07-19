// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppAndInstancePolicyRuleCondition represents the AppAndInstancePolicyRuleCondition schema
// Specifies apps to include or exclude. If `include` is empty, then the condition is met for all apps.
type AppAndInstancePolicyRuleCondition struct {
	// The list of apps or app instances to exclude
	Exclude []interface{} `json:"exclude,omitempty"`
	// The list of apps or app instances to match on
	Include []interface{} `json:"include,omitempty"`
}
