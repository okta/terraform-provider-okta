// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type IdentityProviderApplicationUser struct {
	Embedded    interface{} `json:"_embedded,omitempty"`
	Links       interface{} `json:"_links,omitempty"`
	Created     string      `json:"created,omitempty"`
	ExternalId  string      `json:"externalId,omitempty"`
	Id          string      `json:"id,omitempty"`
	LastUpdated string      `json:"lastUpdated,omitempty"`
	Profile     interface{} `json:"profile,omitempty"`
}
