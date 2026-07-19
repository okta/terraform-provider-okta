// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogStreamSettingsAws represents the LogStreamSettingsAws schema
// Specifies the configuration for the `aws_eventbridge` log stream type. This configuration can't be modified after creation.
type LogStreamSettingsAws struct {
	Region interface{} `json:"region"`
	AccountId interface{} `json:"accountId"`
	EventSourceName interface{} `json:"eventSourceName"`
}
