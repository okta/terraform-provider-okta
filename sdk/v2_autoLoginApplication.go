package sdk

import (
	"time"
)

type AutoLoginApplication struct {
	Embedded      interface{}                   `json:"_embedded,omitempty"`
	Links         interface{}                   `json:"_links,omitempty"`
	Accessibility *ApplicationAccessibility     `json:"accessibility,omitempty"`
	Created       *time.Time                    `json:"created,omitempty"`
	Credentials   *SchemeApplicationCredentials `json:"credentials,omitempty"`
	Features      []string                      `json:"features,omitempty"`
	Id            string                        `json:"id,omitempty"`
	Label         string                        `json:"label,omitempty"`
	LastUpdated   *time.Time                    `json:"lastUpdated,omitempty"`
	Licensing     *ApplicationLicensing         `json:"licensing,omitempty"`
	Name          string                        `json:"name,omitempty"`
	Profile       interface{}                   `json:"profile,omitempty"`
	Settings      *AutoLoginApplicationSettings `json:"settings,omitempty"`
	SignOnMode    string                        `json:"signOnMode,omitempty"`
	Status        string                        `json:"status,omitempty"`
	Visibility    *ApplicationVisibility        `json:"visibility,omitempty"`
}

func NewAutoLoginApplication() *AutoLoginApplication {
	return &AutoLoginApplication{
		SignOnMode: "AUTO_LOGIN",
	}
}

func (a *AutoLoginApplication) IsApplicationInstance() bool {
	return true
}
