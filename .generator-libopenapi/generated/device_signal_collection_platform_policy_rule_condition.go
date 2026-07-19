// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceSignalCollectionPlatformPolicyRuleCondition represents the DeviceSignalCollectionPlatformPolicyRuleCondition schema
// Specifies a particular platform or device to match on
type DeviceSignalCollectionPlatformPolicyRuleCondition struct {
	Exclude []interface{} `json:"exclude,omitempty"`
	Include []interface{} `json:"include,omitempty"`
}
