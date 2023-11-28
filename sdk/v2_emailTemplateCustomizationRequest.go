// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type EmailTemplateCustomizationRequest struct {
	Body      string `json:"body,omitempty"`
	IsDefault *bool  `json:"isDefault,omitempty"`
	Language  string `json:"language,omitempty"`
	Subject   string `json:"subject,omitempty"`
}
