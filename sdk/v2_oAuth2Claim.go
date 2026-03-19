// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type OAuth2Claim struct {
	Links                interface{}            `json:"_links,omitempty"`
	AlwaysIncludeInToken *bool                  `json:"alwaysIncludeInToken,omitempty"`
	ClaimType            string                 `json:"claimType,omitempty"`
	Conditions           *OAuth2ClaimConditions `json:"conditions,omitempty"`
	GroupFilterType      string                 `json:"group_filter_type,omitempty"`
	Id                   string                 `json:"id,omitempty"`
	Name                 string                 `json:"name,omitempty"`
	Status               string                 `json:"status,omitempty"`
	System               *bool                  `json:"system,omitempty"`
	Value                string                 `json:"value,omitempty"`
	ValueType            string                 `json:"valueType,omitempty"`
}
