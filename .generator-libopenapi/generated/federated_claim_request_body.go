// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// FederatedClaimRequestBody represents the FederatedClaimRequestBody schema
type FederatedClaimRequestBody struct {
	// The Okta Expression Language expression to be evaluated at runtime
	Expression string `json:"expression,omitempty"`
	// The name of the claim to be used in the produced token
	Name string `json:"name,omitempty"`
}
