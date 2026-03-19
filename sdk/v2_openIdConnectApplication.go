// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type OpenIdConnectApplication struct {
	Embedded      interface{}                       `json:"_embedded,omitempty"`
	Links         interface{}                       `json:"_links,omitempty"`
	Accessibility *ApplicationAccessibility         `json:"accessibility,omitempty"`
	Created       *time.Time                        `json:"created,omitempty"`
	Credentials   *OAuthApplicationCredentials      `json:"credentials,omitempty"`
	Features      []string                          `json:"features,omitempty"`
	Id            string                            `json:"id,omitempty"`
	Label         string                            `json:"label,omitempty"`
	LastUpdated   *time.Time                        `json:"lastUpdated,omitempty"`
	Licensing     *ApplicationLicensing             `json:"licensing,omitempty"`
	Name          string                            `json:"name,omitempty"`
	Profile       interface{}                       `json:"profile,omitempty"`
	Settings      *OpenIdConnectApplicationSettings `json:"settings,omitempty"`
	SignOnMode    string                            `json:"signOnMode,omitempty"`
	Status        string                            `json:"status,omitempty"`
	Visibility    *ApplicationVisibility            `json:"visibility,omitempty"`
}

func NewOpenIdConnectApplication() *OpenIdConnectApplication {
	return &OpenIdConnectApplication{
		Name:       "oidc_client",
		SignOnMode: "OPENID_CONNECT",
	}
}

func (a *OpenIdConnectApplication) IsApplicationInstance() bool {
	return true
}
