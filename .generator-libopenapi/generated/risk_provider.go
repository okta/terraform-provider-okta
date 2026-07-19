// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// RiskProvider represents the RiskProvider schema
type RiskProvider struct {
	Action interface{} `json:"action"`
	// The ID of the [OAuth 2.0 service app](https://developer.okta.com/docs/guides/implement-oauth-for-okta-serviceapp/main/#create-a-service-app-and-grant-scopes) that's used to send risk events to Okta
	ClientId string `json:"clientId"`
	// Timestamp when the risk provider object was created
	Created *time.Time `json:"created,omitempty"`
	// The ID of the risk provider object
	ID string `json:"id"`
	// Timestamp when the risk provider object was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Name of the risk provider
	Name string `json:"name"`
	Links interface{} `json:"_links"`
}
