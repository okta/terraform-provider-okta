// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PersonalAppsBlockList represents the PersonalAppsBlockList schema
// Defines a list of email domains with a subset of the properties for each domain
type PersonalAppsBlockList struct {
	// List of blocked email domains
	Domains []interface{} `json:"domains,omitempty"`
}
