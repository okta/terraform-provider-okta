package sdk

import (
	"time"
)

type SamlApplication struct {
	Embedded      interface{}               `json:"_embedded,omitempty"`
	Links         interface{}               `json:"_links,omitempty"`
	Accessibility *ApplicationAccessibility `json:"accessibility,omitempty"`
	Created       *time.Time                `json:"created,omitempty"`
	Credentials   *ApplicationCredentials   `json:"credentials,omitempty"`
	Features      []string                  `json:"features,omitempty"`
	Id            string                    `json:"id,omitempty"`
	Label         string                    `json:"label,omitempty"`
	LastUpdated   *time.Time                `json:"lastUpdated,omitempty"`
	Licensing     *ApplicationLicensing     `json:"licensing,omitempty"`
	Name          string                    `json:"name,omitempty"`
	Profile       interface{}               `json:"profile,omitempty"`
	Settings      *SamlApplicationSettings  `json:"settings,omitempty"`
	SignOnMode    string                    `json:"signOnMode,omitempty"`
	Status        string                    `json:"status,omitempty"`
	Visibility    *ApplicationVisibility    `json:"visibility,omitempty"`
}

func NewSamlApplication() *SamlApplication {
	return &SamlApplication{
		SignOnMode: "SAML_2_0",
	}
}

func (a *SamlApplication) IsApplicationInstance() bool {
	return true
}
