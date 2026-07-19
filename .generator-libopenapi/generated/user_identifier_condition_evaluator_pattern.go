// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserIdentifierConditionEvaluatorPattern represents the UserIdentifierConditionEvaluatorPattern schema
// Specifies the details of the patterns to match against
type UserIdentifierConditionEvaluatorPattern struct {
	// The regular expression or simple match string
	Value string `json:"value"`
	MatchType interface{} `json:"matchType"`
}
