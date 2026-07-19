// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RiskDetectionTypesPolicyRuleCondition represents the RiskDetectionTypesPolicyRuleCondition schema
// <x-lifecycle class="oie"></x-lifecycle> An object that references detected risk events. This object can have an `include` parameter or an `exclude` parameter, but not both.
type RiskDetectionTypesPolicyRuleCondition struct {
	// An array of detected risk events to exclude in the entity policy rule
	Exclude []interface{} `json:"exclude,omitempty"`
	// An array of detected risk events to include in the entity policy rule
	Include []interface{} `json:"include,omitempty"`
}
