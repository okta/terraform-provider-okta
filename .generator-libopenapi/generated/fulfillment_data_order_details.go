// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// FulfillmentDataOrderDetails represents the FulfillmentDataOrderDetails schema
// Information about the fulfillment order that includes the factor’s make and model, the custom configuration of the factor, and inventory details.
type FulfillmentDataOrderDetails struct {
	// ID for the set of custom configurations of the requested factor
	CustomizationId string `json:"customizationId,omitempty"`
	// ID for the specific inventory bucket of the requested factor
	InventoryProductId string `json:"inventoryProductId,omitempty"`
	// ID for the make and model of the requested factor
	ProductId string `json:"productId,omitempty"`
}
