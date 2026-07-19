// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogGeolocation represents the LogGeolocation schema
// The latitude and longitude of the geolocation where an action was performed. The object is formatted according to the [ISO 6709](https://www.iso.org/obp/ui/fr/#iso:std:iso:6709:ed-3:v1:en) standard.
type LogGeolocation struct {
	// Latitude which uses two digits for the [integer part](https://www.iso.org/obp/ui/fr/#iso:std:iso:6709:ed-3:v1:en#Latitude)
	Lat float64 `json:"lat,omitempty"`
	// Longitude which uses three digits for the [integer part](https://www.iso.org/obp/ui/fr/#iso:std:iso:6709:ed-3:v1:en#Longitude)
	Lon float64 `json:"lon,omitempty"`
}
