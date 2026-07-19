// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2Claim represents the OAuth2Claim schema
type OAuth2Claim struct {
	GroupFilterType interface{} `json:"group_filter_type,omitempty"`
	// Specifies the value of the Claim. This value must be a string literal if `valueType` is `GROUPS`, and the string literal is matched with the selected `group_filter_type`. The value must be an Okta ...
	Value string `json:"value,omitempty"`
	ValueType interface{} `json:"valueType,omitempty"`
	// Specifies whether to include Claims in the token. The value is always `TRUE` for access token Claims. If the value is set to `FALSE` for an ID token claim, the Claim isn't included in the ID token ...
	AlwaysIncludeInToken bool `json:"alwaysIncludeInToken,omitempty"`
	ClaimType interface{} `json:"claimType,omitempty"`
	// ID of the Claim
	ID string `json:"id,omitempty"`
	// Name of the Claim
	Name string `json:"name,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// When `true`, indicates that Okta created the Claim
	System bool `json:"system,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	Conditions interface{} `json:"conditions,omitempty"`
}
