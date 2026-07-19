// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogDevice represents the LogDevice schema
// The entity that describes a device enrolled with passwordless authentication using Okta Verify.
type LogDevice struct {
	DiskEncryptionType interface{} `json:"disk_encryption_type,omitempty"`
	// Indicates if the device is configured for device management and is registered with Okta
	Managed bool `json:"managed,omitempty"`
	Name string `json:"name,omitempty"`
	OsPlatform string `json:"os_platform,omitempty"`
	// Indicates if the device is registered with an Okta org and is bound to an Okta Verify instance on the device
	Registered bool `json:"registered,omitempty"`
	ScreenLockType interface{} `json:"screen_lock_type,omitempty"`
	// The integration platform or software used with the device
	DeviceIntegrator map[string]interface{} `json:"device_integrator,omitempty"`
	// ID of the device
	ID string `json:"id,omitempty"`
	// If the device has removed software restrictions
	Jailbreak bool `json:"jailbreak,omitempty"`
	OsVersion string `json:"os_version,omitempty"`
	// The availability of hardware security on the device
	SecureHardwarePresent bool `json:"secure_hardware_present,omitempty"`
}
