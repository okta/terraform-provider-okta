// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SalesforceApplicationSettingsApplication represents the SalesforceApplicationSettingsApplication schema
// Salesforce app instance properties
type SalesforceApplicationSettingsApplication struct {
	// Salesforce instance that you want to connect to
	InstanceType string `json:"instanceType"`
	// Salesforce integration type
	IntegrationType string `json:"integrationType"`
	// The Login URL specified in your Salesforce Single Sign-On settings
	LoginUrl string `json:"loginUrl,omitempty"`
	// Salesforce Logout URL
	LogoutUrl string `json:"logoutUrl,omitempty"`
}
