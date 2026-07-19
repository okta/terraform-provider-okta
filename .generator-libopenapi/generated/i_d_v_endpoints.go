// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IDVEndpoints represents the IDVEndpoints schema
// Contains endpoints for the IDV vendor. When you create an `IDV_STANDARD` IdP, you must include the `par`, `authorization`, `token`, and `jwks` endpoints in the request body.
type IDVEndpoints struct {
	Jwks interface{} `json:"jwks"`
	Par interface{} `json:"par"`
	Token interface{} `json:"token"`
	Authorization interface{} `json:"authorization"`
}
