// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EmailTestAddresses represents the EmailTestAddresses schema
type EmailTestAddresses struct {
	// Email address that sends test emails
	FromAddress string `json:"fromAddress"`
	// Email address that receives test emails
	ToAddress string `json:"toAddress"`
}
