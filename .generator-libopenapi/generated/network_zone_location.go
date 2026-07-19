// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// NetworkZoneLocation represents the NetworkZoneLocation schema
type NetworkZoneLocation struct {
	// The two-character ISO 3166-1 country code. Don't use continent codes since they are treated as generic codes for undesignated countries. <br>For example: `US`
	Country string `json:"country,omitempty"`
	// (Optional) The ISO 3166-2 region code appended to the country code (`countryCode-regionCode`), or `null` if empty. Don't use continent codes since they are treated as generic codes for undesignated...
	Region string `json:"region,omitempty"`
}
