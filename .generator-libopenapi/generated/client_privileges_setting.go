// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ClientPrivilegesSetting represents the ClientPrivilegesSetting schema
// The org setting that assigns the super admin role by default to a public client app
type ClientPrivilegesSetting struct {
	// If true, assigns the super admin role by default to new public client apps
	ClientPrivilegesSetting bool `json:"clientPrivilegesSetting,omitempty"`
}
