// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppConfigActiveDirectory represents the AppConfigActiveDirectory schema
type AppConfigActiveDirectory struct {
	GroupType interface{} `json:"groupType"`
	// The SAM account name of the group in Active Directory
	SamAccountName string `json:"samAccountName"`
	// The distinguished name of the group in Active Directory
	DistinguishedName string `json:"distinguishedName"`
	GroupScope interface{} `json:"groupScope"`
}
