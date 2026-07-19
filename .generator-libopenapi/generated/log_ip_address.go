// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogIpAddress represents the LogIpAddress schema
type LogIpAddress struct {
	// Details regarding the source
	Source string `json:"source,omitempty"`
	// IP address version
	Version string `json:"version,omitempty"`
	GeographicalContext interface{} `json:"geographicalContext,omitempty"`
	// IP address
	Ip string `json:"ip,omitempty"`
	IpDetails interface{} `json:"ipDetails,omitempty"`
}
