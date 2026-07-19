// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DevicePostureChecks represents the DevicePostureChecks schema
// <x-lifecycle-container><x-lifecycle class="ea"></x-lifecycle></x-lifecycle-container>Represents the Device Posture Checks configuration for the device assurance policy
type DevicePostureChecks struct {
	// An array of key-value pairs that include the device posture check `variableName` key
	Include []map[string]interface{} `json:"include,omitempty"`
}
