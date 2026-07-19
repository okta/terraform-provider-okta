// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdpDiscoveryPlatformPolicyRuleCondition represents the IdpDiscoveryPlatformPolicyRuleCondition schema
// Specifies a particular platform or device to match on
type IdpDiscoveryPlatformPolicyRuleCondition struct {
	Exclude []interface{} `json:"exclude,omitempty"`
	Include []interface{} `json:"include,omitempty"`
}
