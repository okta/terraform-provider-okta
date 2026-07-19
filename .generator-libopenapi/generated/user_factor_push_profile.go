// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorPushProfile represents the UserFactorPushProfile schema
type UserFactorPushProfile struct {
	// Installed version of Okta Verify
	Version string `json:"version,omitempty"`
	// ID for the factor credential
	CredentialId string `json:"credentialId,omitempty"`
	// Token used to identify the device
	DeviceToken string `json:"deviceToken,omitempty"`
	// Type of device
	DeviceType string `json:"deviceType,omitempty"`
	// Name of the device
	Name string `json:"name,omitempty"`
	// OS version of the associated device
	Platform string `json:"platform,omitempty"`
}
