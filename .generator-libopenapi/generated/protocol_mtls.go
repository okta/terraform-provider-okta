// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProtocolMtls represents the ProtocolMtls schema
// Protocol settings for the [MTLS Protocol](https://tools.ietf.org/html/rfc5246#section-7.4.4)
type ProtocolMtls struct {
	Credentials interface{} `json:"credentials,omitempty"`
	Endpoints interface{} `json:"endpoints,omitempty"`
	// Mutual TLS
	Type string `json:"type,omitempty"`
}
