// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Parameters represents the Parameters schema
// Attributes used for processing Active Directory or LDAP group membership update
type Parameters struct {
	// The update action to take
	Action string `json:"action"`
	// The attribute that tracks group memberships in Active Directory or LDAP. For Active Directory, use `member`. For LDAP, use the appropriate attribute found in the LDAP server such as, but not limite...
	Attribute string `json:"attribute"`
	// List of user IDs whose group memberships to update
	Values []string `json:"values"`
}
