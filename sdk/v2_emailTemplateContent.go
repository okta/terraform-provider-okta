package sdk

type EmailTemplateContent struct {
	Links       interface{} `json:"_links,omitempty"`
	Body        string      `json:"body,omitempty"`
	FromAddress string      `json:"fromAddress,omitempty"`
	FromName    string      `json:"fromName,omitempty"`
	Subject     string      `json:"subject,omitempty"`
}
