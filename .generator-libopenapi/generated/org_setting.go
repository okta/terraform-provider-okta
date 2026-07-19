// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// OrgSetting represents the OrgSetting schema
type OrgSetting struct {
	// Subdomain of org
	Subdomain string `json:"subdomain,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// City of the organization associated with the org
	City string `json:"city,omitempty"`
	// Expiration of org
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Org ID
	ID string `json:"id,omitempty"`
	// State of the organization associated with the org
	State string `json:"state,omitempty"`
	// Status of org
	Status string `json:"status,omitempty"`
	// County of the organization associated with the org
	Country string `json:"country,omitempty"`
	// When org was created
	Created *time.Time `json:"created,omitempty"`
	// Name of org
	CompanyName string `json:"companyName,omitempty"`
	// When org was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Phone number of the organization associated with the org
	PhoneNumber string `json:"phoneNumber,omitempty"`
	// Support help phone of the organization associated with the org
	SupportPhoneNumber string `json:"supportPhoneNumber,omitempty"`
	// Website of the organization associated with the org
	Website string `json:"website,omitempty"`
	// Primary address of the organization associated with the org
	Address1 string `json:"address1,omitempty"`
	// Secondary address of the organization associated with the org
	Address2 string `json:"address2,omitempty"`
	// Support link of org
	EndUserSupportHelpURL string `json:"endUserSupportHelpURL,omitempty"`
	// Postal code of the organization associated with the org
	PostalCode string `json:"postalCode,omitempty"`
}
