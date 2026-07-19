// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProvisioningDetails represents the ProvisioningDetails schema
// Supported provisioning configurations for your integration
type ProvisioningDetails struct {
	Scim interface{} `json:"scim"`
	// List of provisioning features supported in this integration
	Features []string `json:"features"`
}
