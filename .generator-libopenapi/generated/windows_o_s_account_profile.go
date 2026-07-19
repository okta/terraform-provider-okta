// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// WindowsOSAccountProfile represents the WindowsOSAccountProfile schema
type WindowsOSAccountProfile struct {
	// Domain\username format (down-level logon name)
	DownLevelUsername string `json:"downLevelUsername,omitempty"`
	// Full name of the account user
	FullName string `json:"fullName,omitempty"`
	// Globally Unique Identifier for the account
	GUID string `json:"GUID,omitempty"`
	// Fully qualified username
	QualifiedUsername string `json:"qualifiedUsername,omitempty"`
	// Windows Security Identifier (SID)
	SecurityId string `json:"securityId,omitempty"`
	// User principal name
	Upn string `json:"upn,omitempty"`
	// Active Directory join status
	DirectoryJoinStatus string `json:"directoryJoinStatus,omitempty"`
}
