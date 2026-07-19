// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogGeographicalContext represents the LogGeographicalContext schema
// Geographical context describes a set of geographic coordinates. In addition to containing latitude and longitude data, the `GeographicalContext` object also contains address data of postal code-lev...
type LogGeographicalContext struct {
	// The city that encompasses the area that contains the geolocation coordinates, if available (for example, Seattle, San Francisco)
	City string `json:"city,omitempty"`
	// Full name of the country that encompasses the area that contains the geolocation coordinates (for example, France, Uganda)
	Country string `json:"country,omitempty"`
	Geolocation interface{} `json:"geolocation,omitempty"`
	// Postal code of the area that encompasses the geolocation coordinates
	PostalCode string `json:"postalCode,omitempty"`
	// Full name of the state or province that encompasses the area that contains the geolocation coordinates (for example, Montana, Ontario)
	State string `json:"state,omitempty"`
}
