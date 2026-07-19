// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaUserGroupProfile represents the OktaUserGroupProfile schema
// Profile for any group that is not imported from Active Directory. Specifies the standard and custom profile properties for a group.  The `objectClass` for these groups is `okta:user_group`.
type OktaUserGroupProfile struct {
	// Description of the group
	Description string `json:"description,omitempty"`
	// Name of the group
	Name string `json:"name,omitempty"`
}
