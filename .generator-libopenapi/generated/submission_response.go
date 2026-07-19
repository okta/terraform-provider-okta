// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SubmissionResponse represents the SubmissionResponse schema
type SubmissionResponse struct {
	AuthSettings interface{} `json:"authSettings,omitempty"`
	// List of capabilities supported by this integration with embedded protocol configurations
	Capabilities []interface{} `json:"capabilities,omitempty"`
	// A general description of your application and the benefits provided to your customers
	Description string `json:"description,omitempty"`
	GlobalTokenRevocation interface{} `json:"globalTokenRevocation,omitempty"`
	// ID of the user who made the last update
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Provisioning interface{} `json:"provisioning,omitempty"`
	// Status of the OIN Integration submission
	Status string `json:"status,omitempty"`
	// The most recent date and time when a Salesforce case associated with the OIN integration was updated
	CaseLastUpdated string `json:"caseLastUpdated,omitempty"`
	// List of org-level variables for the customer per-tenant configuration. For example, a `subdomain` variable can be used in the ACS URL: `https://${org.subdomain}.example.com/saml/login`
	Config []map[string]interface{} `json:"config,omitempty"`
	// Indicates whether the app submission uses a default logo or a custom-uploaded logo: * If `true`: Uses the default Okta-provided placeholder logo. * If `false`: Uses a custom logo type other than th...
	DefaultLogo bool `json:"defaultLogo,omitempty"`
	// Timestamp when the OIN Integration instance was last updated
	LastUpdated string `json:"lastUpdated,omitempty"`
	// The app integration name. This is the main title used for your integration in the OIN catalog.
	Name string `json:"name,omitempty"`
	// Type of feature supported by the OIN integration
	OinFeatures string `json:"oinFeatures,omitempty"`
	// URL to an uploaded application logo. This logo appears next to your app integration name in the OIN catalog. You must first [Upload an OIN Integration logo](/openapi/okta-management/management/tag/...
	Logo string `json:"logo,omitempty"`
	// List of contact details for the app integration
	AppContactDetails []map[string]interface{} `json:"appContactDetails,omitempty"`
	// OIN Integration ID
	ID string `json:"id,omitempty"`
	// Timestamp when the OIN Integration was last published
	LastPublished string `json:"lastPublished,omitempty"`
	Sso interface{} `json:"sso,omitempty"`
}
