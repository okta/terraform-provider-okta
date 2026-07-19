// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RiskEventSubject represents the RiskEventSubject schema
type RiskEventSubject struct {
	RiskLevel interface{} `json:"riskLevel"`
	// The risk event subject IP address (either an IPv4 or IPv6 address)
	Ip string `json:"ip"`
	// Additional reasons for the risk level of the IP
	Message string `json:"message,omitempty"`
}
