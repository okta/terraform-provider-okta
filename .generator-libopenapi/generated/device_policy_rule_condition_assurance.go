// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DevicePolicyRuleConditionAssurance represents the DevicePolicyRuleConditionAssurance schema
// Specifies [device assurance policies](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/DeviceAssurance/) in the policy rule
type DevicePolicyRuleConditionAssurance struct {
	// Specifies the device assurance policy ID
	Include []string `json:"include,omitempty"`
}
