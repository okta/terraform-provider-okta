// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyDelegationSettingsOptions represents the PasswordPolicyDelegationSettingsOptions schema
type PasswordPolicyDelegationSettingsOptions struct {
	// Indicates if, when performing an unlock operation on an Active Directory sourced User who is locked out of Okta, the system should also attempt to unlock the User's Windows account
	SkipUnlock bool `json:"skipUnlock,omitempty"`
}
