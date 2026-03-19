// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type BookmarkApplication struct {
	Embedded      interface{}                  `json:"_embedded,omitempty"`
	Links         interface{}                  `json:"_links,omitempty"`
	Accessibility *ApplicationAccessibility    `json:"accessibility,omitempty"`
	Created       *time.Time                   `json:"created,omitempty"`
	Credentials   *ApplicationCredentials      `json:"credentials,omitempty"`
	Features      []string                     `json:"features,omitempty"`
	Id            string                       `json:"id,omitempty"`
	Label         string                       `json:"label,omitempty"`
	LastUpdated   *time.Time                   `json:"lastUpdated,omitempty"`
	Licensing     *ApplicationLicensing        `json:"licensing,omitempty"`
	Name          string                       `json:"name,omitempty"`
	Profile       interface{}                  `json:"profile,omitempty"`
	Settings      *BookmarkApplicationSettings `json:"settings,omitempty"`
	SignOnMode    string                       `json:"signOnMode,omitempty"`
	Status        string                       `json:"status,omitempty"`
	Visibility    *ApplicationVisibility       `json:"visibility,omitempty"`
}

func NewBookmarkApplication() *BookmarkApplication {
	return &BookmarkApplication{
		Name:       "bookmark",
		SignOnMode: "BOOKMARK",
	}
}

func (a *BookmarkApplication) IsApplicationInstance() bool {
	return true
}
