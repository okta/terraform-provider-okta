// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicyMappingRequest represents the PolicyMappingRequest schema
type PolicyMappingRequest struct {
	// [Policy ID](https://developer.okta.com/docs/api/openapi/okta-management/management/tags/policy/#tag/Policy/operation/listPolicies!c=200&path=0/id&t=response) of the app sign-in policy that you want...
	ResourceId string `json:"resourceId,omitempty"`
	ResourceType interface{} `json:"resourceType,omitempty"`
}
