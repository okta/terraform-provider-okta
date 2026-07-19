// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AvailableActionProvider represents the AvailableActionProvider schema
type AvailableActionProvider struct {
	// The name of the action flow
	ActionName string `json:"actionName"`
	// The unique identifier of the action flow in the provider system
	ExternalId string `json:"externalId"`
	// The URL to the action flow interface in Workflows platform
	Link string `json:"link"`
	Type interface{} `json:"type"`
}
