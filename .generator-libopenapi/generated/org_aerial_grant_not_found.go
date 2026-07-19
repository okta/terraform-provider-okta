// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OrgAerialGrantNotFound represents the OrgAerialGrantNotFound schema
type OrgAerialGrantNotFound struct {
	// The unique ID of the Aerial account
	AccountId string `json:"accountId,omitempty"`
	// Principal ID of the user who granted the permission
	GrantedBy string `json:"grantedBy,omitempty"`
	// Date when grant was created
	GrantedDate string `json:"grantedDate,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
