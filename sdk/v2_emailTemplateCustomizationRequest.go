package sdk

type EmailTemplateCustomizationRequest struct {
	Body      string `json:"body,omitempty"`
	IsDefault *bool  `json:"isDefault,omitempty"`
	Language  string `json:"language,omitempty"`
	Subject   string `json:"subject,omitempty"`
}
