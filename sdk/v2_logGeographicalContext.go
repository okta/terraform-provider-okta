// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type LogGeographicalContext struct {
	City        string          `json:"city,omitempty"`
	Country     string          `json:"country,omitempty"`
	Geolocation *LogGeolocation `json:"geolocation,omitempty"`
	PostalCode  string          `json:"postalCode,omitempty"`
	State       string          `json:"state,omitempty"`
}
