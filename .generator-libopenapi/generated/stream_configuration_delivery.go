// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// StreamConfigurationDelivery represents the StreamConfigurationDelivery schema
// Contains information about the intended SET delivery method by the receiver
type StreamConfigurationDelivery struct {
	// The HTTP authorization header that's included for each HTTP POST request
	AuthorizationHeader string `json:"authorization_header,omitempty"`
	// The target endpoint URL where the transmitter delivers the SET using HTTP POST requests
	EndpointUrl string `json:"endpoint_url"`
	// The delivery method that the transmitter uses for delivering a SET
	Method string `json:"method"`
}
