// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaSupportCase represents the OktaSupportCase schema
type OktaSupportCase struct {
	// Okta Support case number
	CaseNumber string `json:"caseNumber,omitempty"`
	// Allows the Okta Support team to sign in to your org as an admin and troubleshoot issues
	Impersonation map[string]interface{} `json:"impersonation,omitempty"`
	// Customer allows Okta Support access to self-assigned cases. Support cases are self-assigned when an Okta Support team member creates and assigns the case to themselves.
	SelfAssigned map[string]interface{} `json:"selfAssigned,omitempty"`
	// Subject of the support case
	Subject string `json:"subject,omitempty"`
}
