// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EventHookFilters represents the EventHookFilters schema
// The optional filter defined on a specific event type  > **Note:** Event hook filters is a [self-service Early Access (EA)](/openapi/okta-management/guides/release-lifecycle/#early-access-ea) to ena...
type EventHookFilters struct {
	EventFilterMap interface{} `json:"eventFilterMap,omitempty"`
	// The type of filter. Currently only supports `EXPRESSION_LANGUAGE`
	Type string `json:"type,omitempty"`
}
