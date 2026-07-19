// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ResourceServerJsonWebKeys represents the ResourceServerJsonWebKeys schema
// A [JSON Web Key Set](https://tools.ietf.org/html/rfc7517#section-5) for encrypting JWTs minted by the custom authorization server
type ResourceServerJsonWebKeys struct {
	Keys []interface{} `json:"keys,omitempty"`
}
