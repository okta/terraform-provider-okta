// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationExpressConfiguration represents the ApplicationExpressConfiguration schema
// <div class="x-lifecycle-container"><x-lifecycle class="oie"></x-lifecycle></div> Indicates which Express Configuration capabilities the app supports and has enabled
type ApplicationExpressConfiguration struct {
	// Capabilities currently enabled for the app
	EnabledCapabilities []interface{} `json:"enabledCapabilities,omitempty"`
	// Capabilities supported by the app
	SupportedCapabilities []interface{} `json:"supportedCapabilities,omitempty"`
}
