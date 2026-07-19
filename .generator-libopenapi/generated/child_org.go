// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ChildOrg represents the ChildOrg schema
type ChildOrg struct {
	// Timestamp when the org was created
	Created *time.Time `json:"created,omitempty"`
	// Timestamp when the org was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Unique name of the org. This name appears in the HTML `<title>` tag of the new org sign-in page. Only less than 4-width UTF-8 encoded characters are allowed.
	Name string `json:"name"`
	// Subdomain of the org. Must be unique and include no spaces.
	Subdomain string `json:"subdomain"`
	// Default website for the org
	Website string `json:"website,omitempty"`
	// Specifies available link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json...
	Links map[string]interface{} `json:"_links,omitempty"`
	Admin interface{} `json:"admin"`
	// Edition for the org. `SKU` is the only supported value.
	Edition string `json:"edition"`
	// Org ID
	ID string `json:"id,omitempty"`
	// Settings associated with the created org
	Settings map[string]interface{} `json:"settings,omitempty"`
	// Status of the org. `ACTIVE` is returned after the org is created.
	Status string `json:"status,omitempty"`
	// API token associated with the child org super admin account. Use this API token to provision resources (such as policies, apps, and groups) on the newly created child org. This token is revoked if ...
	Token string `json:"token,omitempty"`
	// Type of returned `token`. See [Okta API tokens](https://developer.okta.com/docs/guides/create-an-api-token/main/#okta-api-tokens).
	TokenType string `json:"tokenType,omitempty"`
}
