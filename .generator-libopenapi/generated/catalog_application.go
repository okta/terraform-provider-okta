// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// CatalogApplication represents the CatalogApplication schema
// An app in the OIN catalog
type CatalogApplication struct {
	// Description of the app in the OIN catalog
	Description string `json:"description,omitempty"`
	// OIN catalog app display name
	DisplayName string `json:"displayName,omitempty"`
	// Features supported by the app. See app [features](/openapi/okta-management/management/application/listapplications#application/listapplications/t=response&c=200&path=&d=0/features).
	Features []string `json:"features,omitempty"`
	// ID of the app instance. Okta returns this property only for apps not in the OIN catalog.
	ID string `json:"id,omitempty"`
	// Timestamp when the object was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// App key name. For OIN catalog apps, this is a unique key for the app definition.
	Name string `json:"name,omitempty"`
	// Authentication mode for the app. See app [signOnMode](/openapi/okta-management/management/application/listapplications#application/listapplications/t=response&c=200&path=&d=0/signonmode).
	SignOnModes []string `json:"signOnModes,omitempty"`
	// Website of the OIN catalog app
	Website string `json:"website,omitempty"`
	// Category for the app in the OIN catalog
	Category string `json:"category,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// OIN verification status of the catalog app
	VerificationStatus string `json:"verificationStatus,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json...
	Links map[string]interface{} `json:"_links,omitempty"`
}
