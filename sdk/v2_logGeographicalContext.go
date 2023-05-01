package sdk

type LogGeographicalContext struct {
	City        string          `json:"city,omitempty"`
	Country     string          `json:"country,omitempty"`
	Geolocation *LogGeolocation `json:"geolocation,omitempty"`
	PostalCode  string          `json:"postalCode,omitempty"`
	State       string          `json:"state,omitempty"`
}
