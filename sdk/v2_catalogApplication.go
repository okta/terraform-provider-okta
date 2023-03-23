package sdk

import (
	"time"
)

type CatalogApplication struct {
	Links              interface{} `json:"_links,omitempty"`
	Category           string      `json:"category,omitempty"`
	Description        string      `json:"description,omitempty"`
	DisplayName        string      `json:"displayName,omitempty"`
	Features           []string    `json:"features,omitempty"`
	Id                 string      `json:"id,omitempty"`
	LastUpdated        *time.Time  `json:"lastUpdated,omitempty"`
	Name               string      `json:"name,omitempty"`
	SignOnModes        []string    `json:"signOnModes,omitempty"`
	Status             string      `json:"status,omitempty"`
	VerificationStatus string      `json:"verificationStatus,omitempty"`
	Website            string      `json:"website,omitempty"`
}
