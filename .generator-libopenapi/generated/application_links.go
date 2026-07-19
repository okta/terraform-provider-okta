// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationLinks represents the ApplicationLinks schema
// Discoverable resources related to the app
type ApplicationLinks struct {
	AccessPolicy interface{} `json:"accessPolicy,omitempty"`
	Activate interface{} `json:"activate,omitempty"`
	Deactivate interface{} `json:"deactivate,omitempty"`
	Groups interface{} `json:"groups,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Self interface{} `json:"self,omitempty"`
	Users interface{} `json:"users,omitempty"`
	// List of app link resources
	AppLinks []interface{} `json:"appLinks,omitempty"`
	Help interface{} `json:"help,omitempty"`
	// List of app logo resources
	Logo []interface{} `json:"logo,omitempty"`
}
