// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogRequest represents the LogRequest schema
// The `Request` object describes details that are related to the HTTP request that triggers this event, if available. When the event isn't sourced to an HTTP request, such as an automatic update on t...
type LogRequest struct {
	// If the incoming request passes through any proxies, the IP addresses of those proxies are stored here in the format of clientIp, proxy1, proxy2, and so on. This field is useful when working with tr...
	IpChain []interface{} `json:"ipChain,omitempty"`
}
