// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DbscScope represents the DbscScope schema
type DbscScope struct {
	// Whether to include subdomains in the binding scope (`false` = exact origin only, `true` = includes subdomains)
	IncludeSite bool `json:"include_site"`
	// The origin URL for which the DBSC binding is valid
	Origin string `json:"origin"`
}
