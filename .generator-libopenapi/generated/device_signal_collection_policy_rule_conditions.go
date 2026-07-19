// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceSignalCollectionPolicyRuleConditions represents the DeviceSignalCollectionPolicyRuleConditions schema
// <x-lifecycle-container><x-lifecycle class="ea"></x-lifecycle></x-lifecycle-container>Specifies conditions that must be met during policy evaluation to apply the rule. All policy conditions, as well...
type DeviceSignalCollectionPolicyRuleConditions struct {
	Network interface{} `json:"network,omitempty"`
	Platform interface{} `json:"platform,omitempty"`
}
