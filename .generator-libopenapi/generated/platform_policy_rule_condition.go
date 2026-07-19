// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PlatformPolicyRuleCondition represents the PlatformPolicyRuleCondition schema
// Specifies a particular platform or device to match on
type PlatformPolicyRuleCondition struct {
	Exclude []interface{} `json:"exclude,omitempty"`
	Include []interface{} `json:"include,omitempty"`
}
