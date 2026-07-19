// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DTCMacOS represents the DTCMacOS schema
// Google Chrome Device Trust Connector provider
type DTCMacOS struct {
	BrowserVersion interface{} `json:"browserVersion,omitempty"`
	// Indicates if a software stack is used to communicate with the DNS server
	BuiltInDnsClientEnabled bool `json:"builtInDnsClientEnabled,omitempty"`
	// Indicates whether access to the Chrome Remote Desktop application is blocked through a policy
	ChromeRemoteDesktopAppBlocked bool `json:"chromeRemoteDesktopAppBlocked,omitempty"`
	// Indicates whether a firewall is enabled at the OS-level on the device
	OsFirewall bool `json:"osFirewall,omitempty"`
	// Indicates whether enterprise-grade (custom) unsafe URL scanning is enabled
	RealtimeUrlCheckMode bool `json:"realtimeUrlCheckMode,omitempty"`
	// Indicates whether the Site Isolation (also known as **Site Per Process**) setting is enabled
	SiteIsolationEnabled bool `json:"siteIsolationEnabled,omitempty"`
	// Enrollment domain of the customer that is currently managing the device
	DeviceEnrollmentDomain string `json:"deviceEnrollmentDomain,omitempty"`
	// Indicates whether the main disk is encrypted
	DiskEncrypted bool `json:"diskEncrypted,omitempty"`
	KeyTrustLevel interface{} `json:"keyTrustLevel,omitempty"`
	OsVersion interface{} `json:"osVersion,omitempty"`
	PasswordProtectionWarningTrigger interface{} `json:"passwordProtectionWarningTrigger,omitempty"`
	SafeBrowsingProtectionLevel interface{} `json:"safeBrowsingProtectionLevel,omitempty"`
	// Indicates whether the device is password-protected
	ScreenLockSecured bool `json:"screenLockSecured,omitempty"`
}
