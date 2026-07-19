// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// FederatedClaim represents the FederatedClaim schema
type FederatedClaim struct {
	// The name of the claim to be used in the produced token
	Name string `json:"name,omitempty"`
	// Timestamp when the federated claim was created
	Created string `json:"created,omitempty"`
	// The Okta Expression Language expression to be evaluated at runtime
	Expression string `json:"expression,omitempty"`
	// The unique ID of the federated claim
	ID string `json:"id,omitempty"`
	// Timestamp when the federated claim was updated
	LastUpdated string `json:"lastUpdated,omitempty"`
}
