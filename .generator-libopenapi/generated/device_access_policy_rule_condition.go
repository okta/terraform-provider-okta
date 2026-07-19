// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceAccessPolicyRuleCondition represents the DeviceAccessPolicyRuleCondition schema
// <x-lifecycle class="oie"></x-lifecycle> Specifies the device condition to match on
type DeviceAccessPolicyRuleCondition struct {
	Assurance interface{} `json:"assurance,omitempty"`
	// Indicates if the device is managed. A device is considered managed if it's part of a device management system.
	Managed bool `json:"managed,omitempty"`
	// Indicates if the device is registered. A device is registered if the User enrolls with Okta Verify that's installed on the device. When the `managed` property is passed, you must also include the `...
	Registered bool `json:"registered,omitempty"`
}
