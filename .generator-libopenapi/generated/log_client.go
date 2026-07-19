// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogClient represents the LogClient schema
// When an event is triggered by an HTTP request, the `client` object describes the [client](https://datatracker.ietf.org/doc/html/rfc2616) that issues the HTTP request. For instance, the web browser ...
type LogClient struct {
	// Type of device that the client operates from (for example, computer)
	Device string `json:"device,omitempty"`
	GeographicalContext interface{} `json:"geographicalContext,omitempty"`
	// For OAuth requests, this is the ID of the OAuth [client](https://datatracker.ietf.org/doc/html/rfc6749#section-1.1) making the request. For SSWS token requests, this is the ID of the agent making t...
	ID string `json:"id,omitempty"`
	// IP address that the client is making its request from
	IpAddress string `json:"ipAddress,omitempty"`
	UserAgent interface{} `json:"userAgent,omitempty"`
	// The `name` of the [Zone](https://developer.okta.com/docs/api/openapi/okta-management/management/networkzone/#tag/NetworkZone/operation/getNetworkZone) that the client's location is mapped to
	Zone string `json:"zone,omitempty"`
}
