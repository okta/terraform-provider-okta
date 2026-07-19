// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogIpDetails represents the LogIpDetails schema
// Details about the associated IP address
type LogIpDetails struct {
	// The [Autonomous system](https://docs.telemetry.mozilla.org/datasets/other/asn_aggregates/reference) number that's associated with the IP address
	AsNumber int `json:"asNumber,omitempty"`
	// The name associated with the Autonomous System Number (ASN)
	AsOrg string `json:"asOrg,omitempty"`
	// The domain name associated with the IP address
	Domain string `json:"domain,omitempty"`
	// The associated IP service categories for the IP address
	IpServiceCategories []interface{} `json:"ipServiceCategories,omitempty"`
	// The internet service provider associated with the IP address
	Isp string `json:"isp,omitempty"`
}
