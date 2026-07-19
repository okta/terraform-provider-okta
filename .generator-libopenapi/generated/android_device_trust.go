// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AndroidDeviceTrust represents the AndroidDeviceTrust schema
// Android Device Trust integration provider
type AndroidDeviceTrust struct {
	DeviceIntegrityLevel interface{} `json:"deviceIntegrityLevel,omitempty"`
	// Indicates whether a device has a network proxy disabled
	NetworkProxyDisabled bool `json:"networkProxyDisabled,omitempty"`
	PlayProtectVerdict interface{} `json:"playProtectVerdict,omitempty"`
	// Indicates whether the device needs to be on the latest major version available to the device  **Note:** This option requires an `osVersion.dynamicVersionRequirement` value to be supplied with the `...
	RequireMajorVersionUpdate bool `json:"requireMajorVersionUpdate,omitempty"`
	ScreenLockComplexity interface{} `json:"screenLockComplexity,omitempty"`
	// Indicates whether Android Debug Bridge (adb) over USB is disabled
	UsbDebuggingDisabled bool `json:"usbDebuggingDisabled,omitempty"`
	// Indicates whether a device is on a password-protected Wi-Fi network
	WifiSecured bool `json:"wifiSecured,omitempty"`
}
