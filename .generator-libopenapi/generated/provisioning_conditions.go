// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProvisioningConditions represents the ProvisioningConditions schema
// Conditional behaviors for an IdP user during authentication
type ProvisioningConditions struct {
	Deprovisioned interface{} `json:"deprovisioned,omitempty"`
	Suspended interface{} `json:"suspended,omitempty"`
}
