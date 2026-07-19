// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DTCWindows represents the DTCWindows schema
// Google Chrome Device Trust Connector provider
type DTCWindows struct {
	// Indicates whether access to the Chrome Remote Desktop application is blocked through a policy
	ChromeRemoteDesktopAppBlocked bool `json:"chromeRemoteDesktopAppBlocked,omitempty"`
	// Customer ID of an installed CrowdStrike agent
	CrowdStrikeCustomerId string `json:"crowdStrikeCustomerId,omitempty"`
	// Enrollment domain of the customer that is currently managing the device
	DeviceEnrollmentDomain string `json:"deviceEnrollmentDomain,omitempty"`
	// Indicates whether enterprise-grade (custom) unsafe URL scanning is enabled
	RealtimeUrlCheckMode bool `json:"realtimeUrlCheckMode,omitempty"`
	// Windows domain that the current machine has joined
	WindowsMachineDomain string `json:"windowsMachineDomain,omitempty"`
	// Windows domain for the current OS user
	WindowsUserDomain string `json:"windowsUserDomain,omitempty"`
	// Indicates whether the main disk is encrypted
	DiskEncrypted bool `json:"diskEncrypted,omitempty"`
	PasswordProtectionWarningTrigger interface{} `json:"passwordProtectionWarningTrigger,omitempty"`
	// Indicates whether the device's startup software has its Secure Boot feature enabled
	SecureBootEnabled bool `json:"secureBootEnabled,omitempty"`
	// Indicates whether the Site Isolation (also known as **Site Per Process**) setting is enabled
	SiteIsolationEnabled bool `json:"siteIsolationEnabled,omitempty"`
	// Indicates whether Chrome is blocking third-party software injection
	ThirdPartyBlockingEnabled bool `json:"thirdPartyBlockingEnabled,omitempty"`
	// Indicates whether antivirus software is enabled
	AntivirusEnabled bool `json:"antivirusEnabled,omitempty"`
	// Agent ID of an installed CrowdStrike agent
	CrowdStrikeAgentId string `json:"crowdStrikeAgentId,omitempty"`
	KeyTrustLevel interface{} `json:"keyTrustLevel,omitempty"`
	// Indicates whether a firewall is enabled at the OS-level on the device
	OsFirewall bool `json:"osFirewall,omitempty"`
	OsVersion interface{} `json:"osVersion,omitempty"`
	SafeBrowsingProtectionLevel interface{} `json:"safeBrowsingProtectionLevel,omitempty"`
	// Indicates whether the device is password-protected
	ScreenLockSecured bool `json:"screenLockSecured,omitempty"`
	BrowserVersion interface{} `json:"browserVersion,omitempty"`
	// Indicates if a software stack is used to communicate with the DNS server
	BuiltInDnsClientEnabled bool `json:"builtInDnsClientEnabled,omitempty"`
}
