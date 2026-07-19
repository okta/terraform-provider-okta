// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceProfile represents the DeviceProfile schema
type DeviceProfile struct {
	// Display name of the device
	DisplayName string `json:"displayName"`
	// Model of the device
	Model string `json:"model,omitempty"`
	// Indicates if the device is registered at Okta
	Registered bool `json:"registered"`
	// Serial number of the device
	SerialNumber string `json:"serialNumber,omitempty"`
	// Windows Trusted Platform Module hash value
	TpmPublicKeyHash string `json:"tpmPublicKeyHash,omitempty"`
	// macOS Unique device identifier of the device
	Udid string `json:"udid,omitempty"`
	DiskEncryptionType interface{} `json:"diskEncryptionType,omitempty"`
	// International Mobile Equipment Identity (IMEI) of the device
	Imei string `json:"imei,omitempty"`
	Platform interface{} `json:"platform"`
	// Indicates if the device contains a secure hardware functionality
	SecureHardwarePresent bool `json:"secureHardwarePresent,omitempty"`
	// Windows Security identifier of the device
	Sid string `json:"sid,omitempty"`
	// Indicates if the device is jailbroken or rooted. Only applicable to `IOS` and `ANDROID` platforms
	IntegrityJailbreak bool `json:"integrityJailbreak,omitempty"`
	// Indicates if the device is managed by mobile device management (MDM) software
	Managed bool `json:"managed,omitempty"`
	// Name of the manufacturer of the device
	Manufacturer string `json:"manufacturer,omitempty"`
	// Mobile equipment identifier of the device
	Meid string `json:"meid,omitempty"`
	// Version of the device OS
	OsVersion string `json:"osVersion,omitempty"`
}
