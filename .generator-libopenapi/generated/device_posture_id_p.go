// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DevicePostureIdP represents the DevicePostureIdP schema
// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Device Posture IdP provider
type DevicePostureIdP struct {
	// Indicates whether the device is compliant according to the custom IDP
	Compliant bool `json:"compliant,omitempty"`
	// Indicates whether the device is managed according to the custom IDP
	Managed bool `json:"managed,omitempty"`
}
