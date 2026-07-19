// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaActiveDirectoryGroupProfile represents the OktaActiveDirectoryGroupProfile schema
// Profile for a group that is imported from Active Directory.  The `objectClass` for such groups is `okta:windows_security_principal`.
type OktaActiveDirectoryGroupProfile struct {
	// Description of the Windows group
	Description string `json:"description,omitempty"`
	// The distinguished name of the Windows group
	Dn string `json:"dn,omitempty"`
	// Base-64 encoded GUID (`objectGUID`) of the Windows group
	ExternalId string `json:"externalId,omitempty"`
	// Name of the Windows group
	Name string `json:"name,omitempty"`
	// Pre-Windows 2000 name of the Windows group
	SamAccountName string `json:"samAccountName,omitempty"`
	// Fully qualified name of the Windows group
	WindowsDomainQualifiedName string `json:"windowsDomainQualifiedName,omitempty"`
}
