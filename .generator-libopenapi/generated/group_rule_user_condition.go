// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupRuleUserCondition represents the GroupRuleUserCondition schema
// Defines conditions specific to user exclusion
type GroupRuleUserCondition struct {
	// Excluded `userIds` when processing rules
	Exclude []string `json:"exclude,omitempty"`
}
