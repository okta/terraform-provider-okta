// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProtocolIdVerification represents the ProtocolIdVerification schema
// Protocol settings for the IDV vendor
type ProtocolIdVerification struct {
	Credentials interface{} `json:"credentials,omitempty"`
	Endpoints interface{} `json:"endpoints,omitempty"`
	Scopes interface{} `json:"scopes,omitempty"`
	// ID verification protocol
	Type string `json:"type,omitempty"`
}
