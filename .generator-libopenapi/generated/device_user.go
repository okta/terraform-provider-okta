// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceUser represents the DeviceUser schema
type DeviceUser struct {
	// The management status of the device
	ManagementStatus string `json:"managementStatus,omitempty"`
	// Screen lock type of the device
	ScreenLockType string `json:"screenLockType,omitempty"`
	User interface{} `json:"user,omitempty"`
	// Timestamp when device was created
	Created string `json:"created,omitempty"`
}
