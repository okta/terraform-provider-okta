// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SubmissionCapability represents the SubmissionCapability schema
// Simple capability structure for capabilities endpoints (PUT/GET /capabilities)
type SubmissionCapability struct {
	Capability interface{} `json:"capability"`
	SupportedProtocols []interface{} `json:"supportedProtocols"`
}
